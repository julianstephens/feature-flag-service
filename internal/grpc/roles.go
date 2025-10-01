package grpcmiddleware

// MethodRoles defines the role requirements for each gRPC method
var MethodRoles = map[string][]string{
	"/FlagService/CreateFlag": {"editor", "admin"},
	"/FlagService/UpdateFlag": {"editor", "admin"},
	"/FlagService/DeleteFlag": {"admin"},
	"/FlagService/GetFlag":    {"user", "editor", "admin"},
	"/FlagService/ListFlags":  {"user", "editor", "admin"},
	// Add user service methods
	"/RbacUserService/CreateUser": {"admin"},
	"/RbacUserService/UpdateUser": {"admin"},
	"/RbacUserService/DeleteUser": {"admin"},
	"/RbacUserService/GetUser":    {"admin"},
	"/RbacUserService/ListUsers":  {"admin"},
}
