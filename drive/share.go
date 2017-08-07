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

    call := self.service.Permissions.Create(args.FileId, permission)

    if permission.Role == "owner" {
        call.TransferOwnership(true);
    }

    _, err := call.Do()

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
	err := self.service.Permissions.Delete(args.FileId, args.PermissionId).Do()
	if err != nil {
		fmt.Errorf("Failed to revoke permission: %s", err)
		return err
	}

	fmt.Fprintf(args.Out, "Permission revoked\n")
	return nil
}

type UpdatePermissionArgs struct {
	Out          io.Writer
	FileId       string
	PermissionId string
	Role         string
}

func (self *Drive) UpdatePermission(args UpdatePermissionArgs) error {
    permission := &drive.Permission{
		Role:               args.Role,
	}

	call := self.service.Permissions.Update(args.FileId, args.PermissionId, permission)

    if permission.Role == "owner" {
        call.TransferOwnership(true);
    }

    _, err := call.Do()

	if err != nil {
		fmt.Errorf("Failed to update permission: %s", err)
		return err
	}

	fmt.Fprintf(args.Out, "Permission updated\n")
	return nil
}

type ListPermissionsArgs struct {
	Out    io.Writer
	SkipHeader  bool
	Separator   string
	FileId string
}

func (self *Drive) ListPermissions(args ListPermissionsArgs) error {
	permList, err := self.service.Permissions.List(args.FileId).Fields("permissions(id,role,type,domain,emailAddress,allowFileDiscovery)").Do()
	if err != nil {
		fmt.Errorf("Failed to list permissions: %s", err)
		return err
	}

	printPermissions(printPermissionsArgs{
		out:         args.Out,
		separator:   args.Separator,
		skipHeader:  args.SkipHeader,
		permissions: permList.Permissions,
	})
	return nil
}

func (self *Drive) shareAnyoneReader(fileId string) error {
	permission := &drive.Permission{
		Role: "reader",
		Type: "anyone",
	}

	_, err := self.service.Permissions.Create(fileId, permission).Do()
	if err != nil {
		return fmt.Errorf("Failed to share file: %s", err)
	}

	return nil
}

type printPermissionsArgs struct {
	out         io.Writer
	skipHeader  bool
	separator   string

	permissions []*drive.Permission
}

func printPermissions(args printPermissionsArgs) {
	w := new(tabwriter.Writer)
	w.Init(args.out, 0, 0, 3, ' ', 0)

	if !args.skipHeader {
		fmt.Fprintf(w, "Id%[1]sType%[1]sRole%[1]sEmail%[1]sDomain%[1]sDiscoverable\n", args.separator)
	}

	for _, p := range args.permissions {
		fmt.Fprintf(w, "%[1]s%[7]s%[2]s%[7]s%[3]s%[7]s%[4]s%[7]s%[5]s%[7]s%[6]s\n",
			p.Id,
			p.Type,
			p.Role,
			p.EmailAddress,
			p.Domain,
			formatBool(p.AllowFileDiscovery),
			args.separator,
		)
	}

	w.Flush()
}
