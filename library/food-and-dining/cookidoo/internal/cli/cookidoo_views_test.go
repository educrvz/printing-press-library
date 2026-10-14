package cli

import "testing"

func TestParseCookidooRecipeRef(t *testing.T) {
	tests := []struct {
		value, locale, id string
	}{
		{value: "r68517", locale: "en-US", id: "r68517"},
		{value: "https://cookidoo.thermomix.com/recipes/recipe/pt-BR/r68517", locale: "pt-BR", id: "r68517"},
	}
	for _, tt := range tests {
		locale, id, _, err := parseCookidooRecipeRef(tt.value, "en-US")
		if err != nil {
			t.Fatalf("parseCookidooRecipeRef(%q): %v", tt.value, err)
		}
		if locale != tt.locale || id != tt.id {
			t.Fatalf("parseCookidooRecipeRef(%q) = locale %q, id %q; want %q, %q", tt.value, locale, id, tt.locale, tt.id)
		}
	}
}

func TestParseCookidooRecipeRefRejectsAccountURLsAndBadIDs(t *testing.T) {
	for _, value := range []string{
		"https://cookidoo.thermomix.com/my-week",
		"https://cookidoo.thermomix.com/recipes/recipe/en-US/not-a-recipe",
		"68517",
	} {
		if _, _, _, err := parseCookidooRecipeRef(value, "en-US"); err == nil {
			t.Fatalf("parseCookidooRecipeRef(%q) unexpectedly succeeded", value)
		}
	}
}

func TestParseCookidooIngredient(t *testing.T) {
	got := parseCookidooIngredient("320 g risotto rice")
	if got.Quantity != 320 || got.Unit != "g" || got.Name != "risotto rice" {
		t.Fatalf("unexpected parsed ingredient: %+v", got)
	}
	got = parseCookidooIngredient("Vegetable stock paste (see Tips)")
	if got.Name != "Vegetable stock paste (see Tips)" || got.Quantity != 0 {
		t.Fatalf("unexpected unquantified ingredient: %+v", got)
	}
}
