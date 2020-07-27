package macro

type Processor func(args []string, scope *Scope) (string, error)
