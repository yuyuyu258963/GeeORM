= range t.NumField() {
		fieldItem := t.Field(i)
		fmt.Println(fieldItem.Name)
		fmt.P