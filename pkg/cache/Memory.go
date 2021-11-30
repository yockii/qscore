package cache

type memoryCacher struct {
	store map[string]string
}

func (c *memoryCacher) GetString(s string) (string, error) {
	if c.store == nil {
		return "", nil
	}
	return c.store[s], nil
}

func (*memoryCacher) Type() string {
	return "memory"
}
