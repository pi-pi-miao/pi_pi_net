package pi_pi_net

func (ctx *Context) Run(network,address string) error {
	ctx, err := ctx.Listen(network, address)
	if err != nil {
		return err
	}
	if err := ctx.Accept();err != nil {
		return err
	}
	return nil
}
