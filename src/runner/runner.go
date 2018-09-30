// this package has a role to pull docker image and run it, then execute commands one by one, finally stop the container
package runner

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"io"
	"os"
	"strings"
	"time"
)

func Run( commit string) {
	config := loadConfig(commit)
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
	}
	pullImage(*cli, ctx, config.Image)
	createdContainer := createContainer(*cli, ctx, config.Image, commit)
	startContainer(createdContainer, ctx, *cli)
	for _, cmd := range config.Steps {
		runCmd(createdContainer, ctx, *cli, strings.Split(cmd.Command, " "))
	}
	displayLog(createdContainer, ctx, *cli)
	var timeOut time.Duration = 1
	cli.ContainerStop(ctx, createdContainer.ID, &timeOut)

}

func pullImage(cli client.Client, ctx context.Context, containerName string) {
	reader, err := cli.ImagePull(ctx, "docker.io/library/"+containerName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)
}

func createContainer(cli client.Client, ctx context.Context, containerName string,commit string) (container.ContainerCreateCreatedBody) {

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	workingDir := currentDir + "/" + commit
	createdContainer, err := cli.ContainerCreate(ctx, &container.Config{
		Image: containerName,
		Tty:   true,
		WorkingDir: "/app",
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: workingDir,
				Target: "/app",
			},
		},
	}, nil, "")
	if err != nil {
		panic(err)
	}
	return createdContainer
}

func startContainer(createdContainer container.ContainerCreateCreatedBody, ctx context.Context, cli client.Client) {
	if err := cli.ContainerStart(ctx, createdContainer.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}

func displayLog(createdContainer container.ContainerCreateCreatedBody, ctx context.Context, cli client.Client) {

	out, err := cli.ContainerLogs(ctx, createdContainer.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)
}

func runCmd(createdContainer container.ContainerCreateCreatedBody, ctx context.Context, cli client.Client, cmd []string) {
	exec, err := cli.ContainerExecCreate(ctx, createdContainer.ID, types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	})
	if err != nil {
		panic(err)
	}
	execResp, err := cli.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{
		Tty: true,
	})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, execResp.Reader)

}
