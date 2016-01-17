package drive

import (
    "fmt"
    "google.golang.org/api/drive/v3"
)


func (self *Drive) Share(args ShareArgs) {
    if args.Revoke {
        err := self.deletePermissions(args)
        errorF(err, "Failed delete permissions: %s", err)
    }

    permission := &drive.Permission{
        AllowFileDiscovery: args.Discoverable,
        Role: args.Role,
        Type: args.Type,
        EmailAddress: args.Email,
    }

    p, err := self.service.Permissions.Create(args.FileId, permission).Do()
    errorF(err, "Failed share file: %s", err)

    fmt.Println(p)
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
