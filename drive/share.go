package drive

import (
    "io"
    "fmt"
    "google.golang.org/api/drive/v3"
)

type ShareArgs struct {
    Out io.Writer
    FileId string
    Role string
    Type string
    Email string
    Discoverable bool
    Revoke bool
}

func (self *Drive) Share(args ShareArgs) (err error) {
    if args.Revoke {
        err = self.deletePermissions(args)
        if err != nil {
            return fmt.Errorf("Failed delete permissions: %s", err)
        }
    }

    permission := &drive.Permission{
        AllowFileDiscovery: args.Discoverable,
        Role: args.Role,
        Type: args.Type,
        EmailAddress: args.Email,
    }

    p, err := self.service.Permissions.Create(args.FileId, permission).Do()
    if err != nil {
        return fmt.Errorf("Failed share file: %s", err)
    }

    fmt.Fprintln(args.Out, p)
    return
}

func (self *Drive) deletePermissions(args ShareArgs) error {
    permList, err := self.service.Permissions.List(args.FileId).Do()
    if err != nil {
        return err
    }

    for _, p := range permList.Permissions {
        // Skip owner permissions
        if p.Role == "owner" {
            continue
        }

        err := self.service.Permissions.Delete(args.FileId, p.Id).Do()
        if err != nil {
            return err
        }
    }

    return nil
}
