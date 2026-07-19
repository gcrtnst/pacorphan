package main

func HelpEnv(t *T) *Env {
	env, err := NewEnv()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := env.Dispose()
		if err != nil {
			t.Log(err)
		}
	})

	env.SetStderr(t.Output())
	return env
}

func HelpMakeAndInstall(t *T, env *Env, src *PkgBuild, explicit bool) {
	err := env.MakeAndInstall(src, explicit)
	if err != nil {
		t.Fatal(err)
	}
}
