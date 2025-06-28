package data

import (
	"fp-kpl/infrastructure/database/schema"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func GetMenus(db *gorm.DB) []schema.Menu {
	var categories []schema.Category

	if err := db.Find(&categories).Error; err != nil {
		log.Fatalf("could not fetch categories: %v", err)
	}
	if len(categories) == 0 {
		log.Fatalf("no categories found")
	}

	categoryMap := make(map[string]string)
	for _, category := range categories {
		categoryMap[category.Name] = category.ID.String()
	}

	return []schema.Menu{
		// Beef
		{
			CategoryID:  getCategoryID(categoryMap, "Beef"),
			Name:        "Beef and Mustard Pie",
			ImageURL:    "https://www.themealdb.com/images/media/meals/sytuqu1511553755.jpg",
			Price:       decimal.NewFromInt(53300),
			IsAvailable: true,
			CookingTime: time.Duration(2 * time.Hour),
			Description: "A beef pie with a mustard sauce",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Beef"),
			Name:        "Beef Brisket Pot Roast",
			ImageURL:    "https://www.themealdb.com/images/media/meals/ursuup1487348423.jpg",
			Price:       decimal.NewFromInt(32000),
			IsAvailable: true,
			CookingTime: time.Duration(4*time.Hour + 15*time.Minute),
			Description: "A fabulous cut of meat. The brisket is located between the shoulders and the forelegs of the steer. These muscles get a workout, which results in more flavor in the meat, and they are also well marbled with fat, adding even more flavor. So they are highly flavorful and perfect for slow braises.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Beef"),
			Name:        "Beef and Oyster Pie",
			ImageURL:    "https://www.themealdb.com/images/media/meals/wrssvt1511556563.jpg",
			Price:       decimal.NewFromInt(55000),
			IsAvailable: true,
			CookingTime: time.Duration(2*time.Hour + 15*time.Minute),
			Description: "Oysters were once as cheap as chips and were used as filled in pies like this. Now a beef and oyster pie is posh enough for a prince. The original surf and turf!",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Beef"),
			Name:        "Beef Bourguignon",
			ImageURL:    "https://www.themealdb.com/images/media/meals/vtqxtu1511784197.jpg",
			Price:       decimal.NewFromInt(58000),
			IsAvailable: true,
			CookingTime: time.Duration(3*time.Hour + 30*time.Minute),
			Description: "A classic French beef stew braised in red wine with onions, garlic, and mushrooms.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Beef"),
			Name:        "Minced Beef Pie",
			ImageURL:    "https://www.themealdb.com/images/media/meals/xwutvy1511555540.jpg",
			Price:       decimal.NewFromInt(46000),
			IsAvailable: true,
			CookingTime: time.Duration(1*time.Hour + 45*time.Minute),
			Description: "A comforting pie filled with savory minced beef and topped with a buttery crust.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Beef"),
			Name:        "Steak and Kidney Pie",
			ImageURL:    "https://www.themealdb.com/images/media/meals/qysyss1511558054.jpg",
			Price:       decimal.NewFromInt(52000),
			IsAvailable: true,
			CookingTime: time.Duration(2*time.Hour + 30*time.Minute),
			Description: "A traditional British pie with slow-cooked steak and kidney in a thick gravy.",
		},

		// Breakfast
		{
			CategoryID:  getCategoryID(categoryMap, "Breakfast"),
			Name:        "Salmon Eggs Eggs Benedict",
			ImageURL:    "https://www.themealdb.com/images/media/meals/1550440197.jpg",
			Price:       decimal.NewFromInt(48000),
			IsAvailable: true,
			CookingTime: time.Duration(45 * time.Minute),
			Description: "A luxurious breakfast classic featuring poached eggs, smoked salmon, and hollandaise on toasted muffins.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Breakfast"),
			Name:        "Home-made Mandazi",
			ImageURL:    "https://www.themealdb.com/images/media/meals/thazgm1555350962.jpg",
			Price:       decimal.NewFromInt(22000),
			IsAvailable: true,
			CookingTime: time.Duration(1 * time.Hour),
			Description: "A popular East African fried bread with a slightly sweet, coconut-infused flavor—perfect for breakfast or tea time.",
		},

		// Chicken
		{
			CategoryID:  getCategoryID(categoryMap, "Chicken"),
			Name:        "Ayam Percik",
			ImageURL:    "https://www.themealdb.com/images/media/meals/020z181619788503.jpg",
			Price:       decimal.NewFromInt(42000),
			IsAvailable: true,
			CookingTime: time.Duration(1*time.Hour + 15*time.Minute),
			Description: "A Malaysian grilled chicken dish marinated in a spicy coconut sauce, bursting with aromatic spices.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Chicken"),
			Name:        "Chicken & Mushroom Hotpot",
			ImageURL:    "https://www.themealdb.com/images/media/meals/uuuspp1511297945.jpg",
			Price:       decimal.NewFromInt(39000),
			IsAvailable: true,
			CookingTime: time.Duration(1*time.Hour + 30*time.Minute),
			Description: "A warm and hearty British hotpot made with chicken, mushrooms, and layered potatoes.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Chicken"),
			Name:        "General Tso's Chicken",
			ImageURL:    "https://www.themealdb.com/images/media/meals/1529444113.jpg",
			Price:       decimal.NewFromInt(47000),
			IsAvailable: true,
			CookingTime: time.Duration(45 * time.Minute),
			Description: "Crispy fried chicken tossed in a sweet, tangy, and spicy sauce—a favorite in Chinese-American cuisine.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Chicken"),
			Name:        "Chicken Alfredo Primavera",
			ImageURL:    "https://www.themealdb.com/images/media/meals/syqypv1486981727.jpg",
			Price:       decimal.NewFromInt(50000),
			IsAvailable: true,
			CookingTime: time.Duration(40 * time.Minute),
			Description: "A creamy pasta dish with sautéed chicken and mixed vegetables in a rich Alfredo sauce.",
		},

		// Dessert
		{
			CategoryID:  getCategoryID(categoryMap, "Dessert"),
			Name:        "Apple & Blackberry Crumble",
			ImageURL:    "https://www.themealdb.com/images/media/meals/xvsurr1511719182.jpg",
			Price:       decimal.NewFromInt(32000),
			IsAvailable: true,
			CookingTime: time.Duration(1 * time.Hour),
			Description: "A warm, crumbly dessert with tart apples and blackberries, topped with golden sugar crust.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Dessert"),
			Name:        "Chocolate Soufflé",
			ImageURL:    "https://www.themealdb.com/images/media/meals/twspvx1511784937.jpg",
			Price:       decimal.NewFromInt(39000),
			IsAvailable: true,
			CookingTime: time.Duration(50 * time.Minute),
			Description: "An airy, melt-in-the-mouth chocolate dessert that's both elegant and indulgent.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Dessert"),
			Name:        "Sticky Toffee Pudding",
			ImageURL:    "https://www.themealdb.com/images/media/meals/xqqqtu1511637379.jpg",
			Price:       decimal.NewFromInt(36000),
			IsAvailable: true,
			CookingTime: time.Duration(1*time.Hour + 10*time.Minute),
			Description: "A moist sponge cake drenched in luscious toffee sauce, served warm with cream or ice cream.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Dessert"),
			Name:        "New York Cheesecake",
			ImageURL:    "https://www.themealdb.com/images/media/meals/swttys1511385853.jpg",
			Price:       decimal.NewFromInt(42000),
			IsAvailable: true,
			CookingTime: time.Duration(1*time.Hour + 30*time.Minute),
			Description: "A classic rich and creamy cheesecake with a buttery graham cracker crust.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Dessert"),
			Name:        "Portuguese Custard Tarts",
			ImageURL:    "https://www.themealdb.com/images/media/meals/vmz7gl1614350221.jpg",
			Price:       decimal.NewFromInt(28000),
			IsAvailable: true,
			CookingTime: time.Duration(1 * time.Hour),
			Description: "Flaky pastry tarts filled with rich egg custard, caramelized on top for a perfect bite.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Dessert"),
			Name:        "Pumpkin Pie",
			ImageURL:    "https://www.themealdb.com/images/media/meals/usuqtp1511385394.jpg",
			Price:       decimal.NewFromInt(35000),
			IsAvailable: true,
			CookingTime: time.Duration(1 * time.Hour),
			Description: "A spiced pumpkin custard pie with a flaky crust, a staple for autumn and festive occasions.",
		},

		// Goat
		{
			CategoryID:  getCategoryID(categoryMap, "Goat"),
			Name:        "Mbuzi Choma (Roasted Goat)",
			ImageURL:    "https://www.themealdb.com/images/media/meals/cuio7s1555492979.jpg",
			Price:       decimal.NewFromInt(58000),
			IsAvailable: true,
			CookingTime: time.Duration(2*time.Hour + 30*time.Minute),
			Description: "A popular East African delicacy made from marinated goat meat, slow-roasted over open flame for a smoky and tender flavor.",
		},

		// Lamb
		{
			CategoryID:  getCategoryID(categoryMap, "Lamb"),
			Name:        "Lamb Rogan Josh",
			ImageURL:    "https://www.themealdb.com/images/media/meals/vvstvq1487342592.jpg",
			Price:       decimal.NewFromInt(57000),
			IsAvailable: true,
			CookingTime: time.Duration(2 * time.Hour),
			Description: "A rich and aromatic Kashmiri lamb curry, slow-cooked with yogurt and warm spices for a bold, flavorful dish.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Lamb"),
			Name:        "Lamb Tzatziki Burgers",
			ImageURL:    "https://www.themealdb.com/images/media/meals/k420tj1585565244.jpg",
			Price:       decimal.NewFromInt(49000),
			IsAvailable: true,
			CookingTime: time.Duration(50 * time.Minute),
			Description: "Juicy lamb burgers seasoned with Mediterranean spices and topped with cool, creamy tzatziki in a soft bun.",
		},

		// Miscellaneous
		{
			CategoryID:  getCategoryID(categoryMap, "Miscellaneous"),
			Name:        "Duck Confit",
			ImageURL:    "https://www.themealdb.com/images/media/meals/wvpvsu1511786158.jpg",
			Price:       decimal.NewFromInt(72000),
			IsAvailable: true,
			CookingTime: time.Duration(3 * time.Hour),
			Description: "A French classic—duck legs slowly cooked in their own fat until meltingly tender, then crisped before serving.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Miscellaneous"),
			Name:        "Chakchouka",
			ImageURL:    "https://www.themealdb.com/images/media/meals/gpz67p1560458984.jpg",
			Price:       decimal.NewFromInt(30000),
			IsAvailable: true,
			CookingTime: time.Duration(40 * time.Minute),
			Description: "A North African dish of poached eggs in a spicy tomato and pepper sauce—simple, hearty, and flavorful.",
		},

		// Pasta
		{
			CategoryID:  getCategoryID(categoryMap, "Pasta"),
			Name:        "Chilli Prawn Linguine",
			ImageURL:    "https://www.themealdb.com/images/media/meals/usywpp1511189717.jpg",
			Price:       decimal.NewFromInt(55000),
			IsAvailable: true,
			CookingTime: time.Duration(45 * time.Minute),
			Description: "A spicy and zesty pasta dish with succulent prawns, garlic, chili, and linguine noodles.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Pasta"),
			Name:        "Lasagne",
			ImageURL:    "https://www.themealdb.com/images/media/meals/wtsvxx1511296896.jpg",
			Price:       decimal.NewFromInt(52000),
			IsAvailable: true,
			CookingTime: time.Duration(1*time.Hour + 30*time.Minute),
			Description: "A traditional Italian layered pasta dish with rich meat sauce, creamy béchamel, and melted cheese.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Pasta"),
			Name:        "Spaghetti alla Carbonara",
			ImageURL:    "https://www.themealdb.com/images/media/meals/llcbn01574260722.jpg",
			Price:       decimal.NewFromInt(48000),
			IsAvailable: true,
			CookingTime: time.Duration(30 * time.Minute),
			Description: "An Italian classic with creamy egg sauce, crispy pancetta, and freshly grated cheese over spaghetti.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Pasta"),
			Name:        "Grilled Mac and Cheese Sandwich",
			ImageURL:    "https://www.themealdb.com/images/media/meals/xutquv1505330523.jpg",
			Price:       decimal.NewFromInt(45000),
			IsAvailable: true,
			CookingTime: time.Duration(25 * time.Minute),
			Description: "Comfort food to the max—gooey mac and cheese grilled between buttery slices of toasted bread.",
		},

		// Pork
		{
			CategoryID:  getCategoryID(categoryMap, "Pork"),
			Name:        "Bigos (Hunters Stew)",
			ImageURL:    "https://www.themealdb.com/images/media/meals/md8w601593348504.jpg",
			Price:       decimal.NewFromInt(51000),
			IsAvailable: true,
			CookingTime: time.Duration(2*time.Hour + 15*time.Minute),
			Description: "A traditional Polish stew made with various cuts of pork, sauerkraut, and spices—hearty and full of flavor.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Pork"),
			Name:        "Japanese Katsudon",
			ImageURL:    "https://www.themealdb.com/images/media/meals/d8f6qx1604182128.jpg",
			Price:       decimal.NewFromInt(49000),
			IsAvailable: true,
			CookingTime: time.Duration(50 * time.Minute),
			Description: "A comforting Japanese rice bowl topped with breaded pork cutlet, egg, and a savory-sweet soy sauce broth.",
		},

		// Seafood
		{
			CategoryID:  getCategoryID(categoryMap, "Seafood"),
			Name:        "Fish Pie",
			ImageURL:    "https://www.themealdb.com/images/media/meals/ysxwuq1487323065.jpg",
			Price:       decimal.NewFromInt(48000),
			IsAvailable: true,
			CookingTime: time.Duration(1*time.Hour + 15*time.Minute),
			Description: "A creamy baked dish made with white fish, smoked fish, and shrimp, topped with mashed potatoes.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Seafood"),
			Name:        "Escovitch Fish",
			ImageURL:    "https://www.themealdb.com/images/media/meals/1520084413.jpg",
			Price:       decimal.NewFromInt(47000),
			IsAvailable: true,
			CookingTime: time.Duration(1 * time.Hour),
			Description: "A Jamaican dish of crispy fried fish topped with spicy pickled vegetables and vinegar sauce.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Seafood"),
			Name:        "Honey Teriyaki Salmon",
			ImageURL:    "https://www.themealdb.com/images/media/meals/xxyupu1468262513.jpg",
			Price:       decimal.NewFromInt(52000),
			IsAvailable: true,
			CookingTime: time.Duration(40 * time.Minute),
			Description: "Tender salmon fillets glazed with a sweet and savory honey teriyaki sauce.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Seafood"),
			Name:        "Laksa King Prawn Noodles",
			ImageURL:    "https://www.themealdb.com/images/media/meals/rvypwy1503069308.jpg",
			Price:       decimal.NewFromInt(55000),
			IsAvailable: true,
			CookingTime: time.Duration(1 * time.Hour),
			Description: "A rich, spicy noodle soup with prawns, coconut milk, and Southeast Asian herbs.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Seafood"),
			Name:        "Recheado Masala Fish",
			ImageURL:    "https://www.themealdb.com/images/media/meals/uwxusv1487344500.jpg",
			Price:       decimal.NewFromInt(46000),
			IsAvailable: true,
			CookingTime: time.Duration(1 * time.Hour),
			Description: "Goan-style stuffed fish with tangy, spicy recheado masala paste and pan-fried to perfection.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Seafood"),
			Name:        "Sushi",
			ImageURL:    "https://www.themealdb.com/images/media/meals/g046bb1663960946.jpg",
			Price:       decimal.NewFromInt(60000),
			IsAvailable: true,
			CookingTime: time.Duration(1*time.Hour + 30*time.Minute),
			Description: "A Japanese delicacy of vinegared rice combined with raw or cooked seafood and vegetables, beautifully rolled and sliced.",
		},

		// Side
		{
			CategoryID:  getCategoryID(categoryMap, "Side"),
			Name:        "Burek",
			ImageURL:    "https://www.themealdb.com/images/media/meals/tkxquw1628771028.jpg",
			Price:       decimal.NewFromInt(35000),
			IsAvailable: true,
			CookingTime: time.Duration(1 * time.Hour),
			Description: "A flaky, savory pastry filled with spiced minced meat, popular across the Balkans and Turkey.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Side"),
			Name:        "French Onion Soup",
			ImageURL:    "https://www.themealdb.com/images/media/meals/xvrrux1511783685.jpg",
			Price:       decimal.NewFromInt(39000),
			IsAvailable: true,
			CookingTime: time.Duration(1*time.Hour + 10*time.Minute),
			Description: "A rich and hearty soup made with caramelized onions and beef stock, topped with toasted bread and melted cheese.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Side"),
			Name:        "Pierogi (Polish Dumplings)",
			ImageURL:    "https://www.themealdb.com/images/media/meals/45xxr21593348847.jpg",
			Price:       decimal.NewFromInt(37000),
			IsAvailable: true,
			CookingTime: time.Duration(1*time.Hour + 20*time.Minute),
			Description: "Soft dumplings stuffed with potato, cheese, or meat, traditionally boiled and then pan-fried in butter.",
		},

		// Starter
		{
			CategoryID:  getCategoryID(categoryMap, "Starter"),
			Name:        "Clam Chowder",
			ImageURL:    "https://www.themealdb.com/images/media/meals/rvtvuw1511190488.jpg",
			Price:       decimal.NewFromInt(45000),
			IsAvailable: true,
			CookingTime: time.Duration(1 * time.Hour),
			Description: "A creamy, hearty soup made with clams, potatoes, and celery — a comforting New England classic.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Starter"),
			Name:        "Cream Cheese Tart",
			ImageURL:    "https://www.themealdb.com/images/media/meals/wurrux1468416624.jpg",
			Price:       decimal.NewFromInt(38000),
			IsAvailable: true,
			CookingTime: time.Duration(1*time.Hour + 15*time.Minute),
			Description: "A light and creamy tart with a buttery crust, filled with smooth cream cheese — perfect for dessert or tea time.",
		},

		// Vegan
		{
			CategoryID:  getCategoryID(categoryMap, "Vegan"),
			Name:        "Roast Fennel and Aubergine Paella",
			ImageURL:    "https://www.themealdb.com/images/media/meals/1520081754.jpg",
			Price:       decimal.NewFromInt(47000),
			IsAvailable: true,
			CookingTime: time.Duration(1 * time.Hour),
			Description: "A vibrant vegan twist on classic Spanish paella with roasted fennel, eggplant, and bold spices.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Vegan"),
			Name:        "Vegan Chocolate Cake",
			ImageURL:    "https://www.themealdb.com/images/media/meals/qxutws1486978099.jpg",
			Price:       decimal.NewFromInt(42000),
			IsAvailable: true,
			CookingTime: time.Duration(1*time.Hour + 10*time.Minute),
			Description: "Rich and moist dairy-free chocolate cake made with plant-based ingredients, perfect for dessert lovers.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Vegan"),
			Name:        "Vegan Lasagna",
			ImageURL:    "https://www.themealdb.com/images/media/meals/rvxxuy1468312893.jpg",
			Price:       decimal.NewFromInt(50000),
			IsAvailable: true,
			CookingTime: time.Duration(1*time.Hour + 30*time.Minute),
			Description: "Layered pasta with rich tomato sauce, creamy vegan béchamel, and vegetable filling — hearty and comforting.",
		},

		// Vegetarian
		{
			CategoryID:  getCategoryID(categoryMap, "Vegetarian"),
			Name:        "Baingan Bharta",
			ImageURL:    "https://www.themealdb.com/images/media/meals/urtpqw1487341253.jpg",
			Price:       decimal.NewFromInt(42000),
			IsAvailable: true,
			CookingTime: time.Duration(1 * time.Hour),
			Description: "A smoky North Indian dish made from roasted and mashed eggplant cooked with onions, tomatoes, and spices.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Vegetarian"),
			Name:        "Cabbage Soup (Shchi)",
			ImageURL:    "https://www.themealdb.com/images/media/meals/60oc3k1699009846.jpg",
			Price:       decimal.NewFromInt(37000),
			IsAvailable: true,
			CookingTime: time.Duration(1*time.Hour + 15*time.Minute),
			Description: "A traditional Russian soup made with cabbage, potatoes, and herbs — light, healthy, and warming.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Vegetarian"),
			Name:        "Matar Paneer",
			ImageURL:    "https://www.themealdb.com/images/media/meals/xxpqsy1511452222.jpg",
			Price:       decimal.NewFromInt(46000),
			IsAvailable: true,
			CookingTime: time.Duration(50 * time.Minute),
			Description: "A creamy North Indian curry made with green peas and paneer cheese in a spiced tomato gravy.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Vegetarian"),
			Name:        "Koshari",
			ImageURL:    "https://www.themealdb.com/images/media/meals/4er7mj1598733193.jpg",
			Price:       decimal.NewFromInt(40000),
			IsAvailable: true,
			CookingTime: time.Duration(1 * time.Hour),
			Description: "A popular Egyptian street food of rice, lentils, pasta, and chickpeas topped with spicy tomato sauce and crispy onions.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Vegetarian"),
			Name:        "Stuffed Bell Peppers with Quinoa and Black Beans",
			ImageURL:    "https://www.themealdb.com/images/media/meals/b66myb1683207208.jpg",
			Price:       decimal.NewFromInt(48000),
			IsAvailable: true,
			CookingTime: time.Duration(1*time.Hour + 10*time.Minute),
			Description: "Bell peppers stuffed with a hearty mix of quinoa, black beans, corn, and spices — healthy and satisfying.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Vegetarian"),
			Name:        "Yaki Udon",
			ImageURL:    "https://www.themealdb.com/images/media/meals/wrustq1511475474.jpg",
			Price:       decimal.NewFromInt(43000),
			IsAvailable: true,
			CookingTime: time.Duration(40 * time.Minute),
			Description: "Japanese stir-fried udon noodles with mixed vegetables and soy-based sauce — savory and umami-rich.",
		},
	}
}

func getCategoryID(m map[string]string, name string) uuid.UUID {
	idStr, ok := m[name]
	if !ok {
		log.Fatalf("category %s not found", name)
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Fatalf("invalid category ID: %v", err)
	}

	return id
}
