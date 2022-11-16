package main

import (
    "context"
    "fmt"
    "os"

    "dagger.io/dagger"
)

func main() {
    if err := build(context.Background()); err != nil {
        fmt.Println(err)
    }
}

func build(ctx context.Context) error {
    fmt.Println("Building with Dagger")

    // daggerクライアントの初期化
    client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
    if err != nil {
        return err
    }
    defer client.Close()

    // カレントディレクトリのパスを取得
    src := client.Host().Workdir()

    // `golang` コンテナイメージを持ってくる
    golang := client.Container().From("golang:latest")

    // `golang` コンテナにカレントディレクトリをマウントする
    golang = golang.WithMountedDirectory("/src", src).WithWorkdir("/")

    // `golang` コンテナでrootをlsする
    golang = golang.Exec(dagger.ContainerExecOpts{
        Args: []string{"ls", "/"},
    })

    ls, err := golang.Stdout().Contents(ctx)
    if err != nil {
        return err
    }
    fmt.Println(ls)

    // `golang` コンテナでhello world!
    golang = golang.Exec(dagger.ContainerExecOpts{
        Args: []string{"echo", "'hello world!'"},
    })

    echo, err := golang.Stdout().Contents(ctx)
    if err != nil {
        return err
    }
    fmt.Println(echo)

    return nil
}
