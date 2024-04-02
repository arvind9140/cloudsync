package cmd

import (
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "github.com/spf13/cobra"
)
var sourceDir string
var destDir string


var syncCmdLocalToLocal = &cobra.Command{
    Use:   "sync_local",
    Short: "Sync files and directories",
    Run: func(cmd *cobra.Command, args []string) {
        err := syncDirectories(sourceDir, destDir)
        if err != nil {
            fmt.Println("Error syncing:", err)
            return
        }
        fmt.Println("Sync completed successfully!")
    },
}

func init() {
    syncCmdLocalToLocal.Flags().StringVarP(&sourceDir, "source", "s", "", "Source directory")
    syncCmdLocalToLocal.Flags().StringVarP(&destDir, "destination", "d", "", "Destination directory")

    syncCmdLocalToLocal.MarkFlagRequired("source")
    syncCmdLocalToLocal.MarkFlagRequired("destination")

    rootCmd.AddCommand(syncCmdLocalToLocal)
}


func syncDirectories(sourceDir, destDir string) error {
    sourceFiles, err := ioutil.ReadDir(sourceDir)
    if err != nil {
        return err
    }

    for _, file := range sourceFiles {
        sourceFilePath := filepath.Join(sourceDir, file.Name())
        destFilePath := filepath.Join(destDir, file.Name())

        if file.IsDir() {
            err := syncDirectories(sourceFilePath, destFilePath)
            if err != nil {
                return err
            }
        } else {
            err := copyFile(sourceFilePath, destFilePath)
            if err != nil {
                return err
            }
        }
    }

    return nil
}

func copyFile(src, dst string) error {
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer dstFile.Close()

    _, err = io.Copy(dstFile, srcFile)
    if err != nil {
        return err
    }

    return nil
}
