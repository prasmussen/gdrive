package drive

import (
	"fmt"
	"google.golang.org/api/drive/v3"
	"io"
	"text/tabwriter"
)

type ShareArgs struct {
	Out          io.Writer
	FileId       string
	Role         string
	Type         string
	Email        string
	Domain       string
	Discoverable bool
}

func (self *Drive) Share(args ShareArgs) error {
	permission := &drive.Permission{
		AllowFileDiscovery: args.Discoverable,
		Role:               args.Role,
		Type:               args.Type,
		EmailAddress:       args.Email,
		Domain:             args.Domain,
	}

	_, err := self.service.Permissions.Create(args.FileId, permission).SupportsTeamDrives(true).Do()
	if err != nil {
		return fmt.Errorf("Failed to share file: %s", err)
	}

	fmt.Fprintf(args.Out, "Granted %s permission to %s\n", args.Role, args.Type)
	return nil
}

type RevokePermissionArgs struct {
	Out          io.Writer
	FileId       string
	PermissionId string
}

func (self *Drive) RevokePermission(args RevokePermissionArgs) error {
	err := self.service.Permissions.Delete(args.FileId, args.PermissionId).SupportsTeamDrives(true).Do()
	if err != nil {
		fmt.Errorf("Failed to revoke permission: %s", err)
		return err
	}

	fmt.Fprintf(args.Out, "Permission revoked\n")
	return nil
}

type ListPermissionsArgs struct {
	Out    io.Writer
	FileId string
}

func (self *Drive) ListPermissions(args ListPermissionsArgs) error {
	permList, err := self.service.Permissions.List(args.FileId).SupportsTeamDrives(true).Fields("permissions(id,role,type,domain,emailAddress,allowFileDiscovery)").Do()
	if err != nil {
		fmt.Errorf("Failed to list permissions: %s", err)
		return err
	}

	printPermissions(printPermissionsArgs{
		out:         args.Out,
		permissions: permList.Permissions,
	})
	return nil
}

func (self *Drive) shareAnyoneReader(fileId string) error {
	permission := &drive.Permission{
		Role: "reader",
		Type: "anyone",
	}

	_, err := self.service.Permissions.Create(fileId, permission).SupportsTeamDrives(true).Do()
	if err != nil {
		return fmt.Errorf("Failed to share file: %s", err)
	}

	return nil
}

type printPermissionsArgs struct {
	out         io.Writer
	permissions []*drive.Permission
}

func printPermissions(args printPermissionsArgs) {
	w := new(tabwriter.Writer)
	w.Init(args.out, 0, 0, 3, ' ', 0)

	fmt.Fprintln(w, "Id\tType\tRole\tEmail\tDomain\tDiscoverable")

	for _, p := range args.permissions {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			p.Id,
			p.Type,
			p.Role,
			p.EmailAddress,
			p.Domain,
			formatBool(p.AllowFileDiscovery),
		)
	}

	w.Flush()
}
