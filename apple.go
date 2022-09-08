package main

import (
	"fmt"
	"strings"
)

// appleModels is different Apple models and their known IDs.
// Sources:
// https://www.theiphonewiki.com/wiki/Models
// https://everymac.com/systems/by_capability/mac-specs-by-machine-model-machine-id.html
var appleModels = map[string][]string{
	"Apple TV (1st generation)":    nil, // unknown id
	"Apple TV (2nd generation)":    {"K66AP"},
	"Apple TV (3rd generation)":    {"J33AP", "J33IAP"},
	"Apple TV (4th generation)":    {"J42dAP"},
	"Apple TV 4K":                  {"J105aAP"},
	"Apple TV 4K (2nd generation)": {"J305AP"},

	"iPad":                  {"K48AP"},
	"iPad 2":                {"K93AP", "K94AP", "K95AP", "K93AAP"},
	"iPad (3rd generation)": {"J1AP", "J2AP", "J2AAP"},
	"iPad (4th generation)": {"P101AP", "P102AP", "P103AP"},
	"iPad (5th generation)": {"J71sAP", "J71tAP", "J72sAP", "J72tAP"},
	"iPad (6th generation)": {"J71bAP", "J72bAP"},
	"iPad (7th generation)": {"J171AP", "J172AP"},
	"iPad (8th generation)": {"J171aAP", "J172aAP"},
	"iPad (9th generation)": {"J181AP", "J182AP"},

	"iPad Air":                  {"J71AP", "J72AP", "J73AP"},
	"iPad Air 2":                {"J81AP", "J82AP"},
	"iPad Air (3rd generation)": {"J217AP", "J218AP"},
	"iPad Air (4th generation)": {"J307AP", "J308AP"},
	"iPad Air (5th generation)": {"J407AP", "J408AP"},

	"iPad Pro (12.9-inch)":                  {"J98aAP", "J99aAP"},
	"iPad Pro (9.7-inch)":                   {"J127AP", "J128AP"},
	"iPad Pro (12.9-inch) (2nd generation)": {"J120AP", "J121AP"},
	"iPad Pro (10.5-inch)":                  {"J207AP", "J208AP"},
	"iPad Pro (11-inch)":                    {"J317AP", "J317xAP", "J318AP", "J318xAP"},
	"iPad Pro (12.9-inch) (3rd generation)": {"J320AP", "J320xAP", "J321AP", "J321xAP"},
	"iPad Pro (11-inch) (2nd generation)":   {"J417AP", "J418AP"},
	"iPad Pro (12.9-inch) (4th generation)": {"J420AP", "J421AP"},
	"iPad Pro (11-inch) (3rd generation)":   {"J517AP", "J517xAP", "J518AP", "J518xAP"},
	"iPad Pro (12.9-inch) (5th generation)": {"J522AP", "J522xAP", "J523AP", "J523xAP"},

	"iPad mini":                  {"P105AP", "P106AP", "P107AP"},
	"iPad mini 2":                {"J85AP", "J86AP", "J87AP"},
	"iPad mini 3":                {"J85mAP", "J86mAP", "J87mAP"},
	"iPad mini 4":                {"J96AP", "J97AP"},
	"iPad mini (5th generation)": {"J210AP", "J211AP"},
	"iPad mini (6th generation)": {"J310AP", "J311AP"},

	"iPhone":                     {"M68AP"},
	"iPhone 3G":                  {"N82AP"},
	"iPhone 3GS":                 {"N88AP"},
	"iPhone 4":                   {"N90AP", "N90bAP", "N92AP"},
	"iPhone 4S":                  {"N94AP"},
	"iPhone 5":                   {"N41AP", "N42AP"},
	"iPhone 5c":                  {"N48AP", "N49AP"},
	"iPhone 5s":                  {"N51AP", "N53AP"},
	"iPhone 6":                   {"N61AP"},
	"iPhone 6 Plus":              {"N56AP"},
	"iPhone 6s":                  {"N71AP", "N71mAP"},
	"iPhone 6s Plus":             {"N66AP", "N66mAP"},
	"iPhone SE (1st generation)": {"N69AP", "N69uAP"},
	"iPhone 7":                   {"D10AP", "D101AP"},
	"iPhone 7 Plus":              {"D11AP", "D111AP"},
	"iPhone 8":                   {"D20AP", "D20AAP", "D201AP", "D201AAP"},
	"iPhone 8 Plus":              {"D21AP", "D21AAP", "D211AP", "D211AAP"},
	"iPhone X":                   {"D22AP", "D221AP"},
	"iPhone XR":                  {"N841AP"},
	"iPhone XS":                  {"D321AP"},
	"iPhone XS Max":              {"D331pAP"},
	"iPhone 11":                  {"N104AP"},
	"iPhone 11 Pro":              {"D421AP"},
	"iPhone 11 Pro Max":          {"D431AP"},
	"iPhone SE (2nd generation)": {"D79AP"},
	"iPhone 12 mini":             {"D52gAP"},
	"iPhone 12":                  {"D53gAP"},
	"iPhone 12 Pro":              {"D53pAP"},
	"iPhone 12 Pro Max":          {"D54pAP"},
	"iPhone 13 Mini":             {"D16AP"},
	"iPhone 13":                  {"D17AP"},
	"iPhone 13 Pro":              {"D63AP"},
	"iPhone 13 Pro Max":          {"D64AP"},
	"iPhone SE (3rd generation)": {"D49AP"},

	"eMac G4":                       {"PowerMac4,4", "PowerMac6,4"},
	"iBook G3":                      {"PowerBook2,1", "PowerBook4,1", "PowerBook4,2", "PowerBook4,3"},
	"iBook G4":                      {"PowerBook6,3", "PowerBook6,5", "PowerBook6,7"},
	"iMac 17-inch":                  {"iMac4,2", "iMac5,2"},
	"iMac 17/20-inch":               {"iMac4,1", "iMac5,1"},
	"iMac 20/24-inch":               {"iMac7,1", "iMac8,1", "iMac9,1"},
	"iMac 21.5-inch":                {"iMac11,2", "iMac12,1", "iMac13,1", "iMac14,1", "iMac14,3", "iMac14,4", "iMac16,1", "iMac16,2", "iMac18,1", "iMac18,2", "iMac19,2"},
	"iMac 21.5/27-inch":             {"iMac10,1"},
	"iMac 24-inch":                  {"iMac6,1"},
	"iMac 27-inch":                  {"iMac11,1", "iMac11,3", "iMac12,2", "iMac13,2", "iMac14,2", "iMac15,1", "iMac17,1", "iMac18,3", "iMac19,1"},
	"iMac G3":                       {"iMac,1", "PowerMac2,1", "PowerMac2,2", "PowerMac4,1"},
	"iMac G4":                       {"PowerMac4,2", "PowerMac4,5", "PowerMac6,1", "PowerMac6,3"},
	"iMac G5":                       {"PowerMac8,1", "PowerMac8,2", "PowerMac12,1"},
	"iMac Pro":                      {"iMacPro1,1"},
	"Mac Mini G4":                   {"PowerMac10,1", "PowerMac10,2"},
	"Mac Mini Intel":                {"Macmini1,1", "Macmini2,1", "Macmini3,1", "Macmini4,1", "Macmini5,1", "Macmini5,2", "Macmini5,3", "Macmini6,1", "Macmini6,2", "Macmini7,1", "Macmini8,1"},
	"Mac Pro":                       {"MacPro1,1*", "MacPro1,1", "MacPro2,1", "MacPro3,1", "MacPro4,1", "MacPro5,1", "MacPro6,1", "MacPro7,1"},
	"MacBook 12-inch":               {"MacBook8,1", "MacBook9,1", "MacBook10,1"},
	"MacBook 13-inch":               {"MacBook1,1", "MacBook2,1", "MacBook3,1", "MacBook4,1", "MacBook5,1", "MacBook5,2", "MacBook6,1", "MacBook7,1"},
	"MacBook Air 11-inch":           {"MacBookAir3,1", "MacBookAir4,1", "MacBookAir5,1", "MacBookAir6,1", "MacBookAir7,1"},
	"MacBook Air 13-inch":           {"MacBookAir1,1", "MacBookAir2,1", "MacBookAir3,2", "MacBookAir4,2", "MacBookAir5,2", "MacBookAir6,2", "MacBookAir7,2", "MacBookAir8,1", "MacBookAir8,2", "MacBookAir9,1"},
	"MacBook Pro 13-inch":           {"MacBookPro5,5", "MacBookPro7,1", "MacBookPro8,1", "MacBookPro9,2", "MacBookPro10,2", "MacBookPro11,1", "MacBookPro12,1", "MacBookPro13,1", "MacBookPro14,1", "MacBookPro15,2", "MacBookPro15,4", "MacBookPro16,2", "MacBookPro16,4"},
	"MacBook Pro 13-inch Touch":     {"MacBookPro13,2", "MacBookPro14,2"},
	"MacBook Pro 15-inch":           {"MacBookPro1,1", "MacBookPro2,2", "MacBookPro5,1", "MacBookPro5,3", "MacBookPro5,4", "MacBookPro6,2", "MacBookPro8,2", "MacBookPro9,1", "MacBookPro10,1", "MacBookPro11,2", "MacBookPro11,3", "MacBookPro11,4", "MacBookPro11,5"},
	"MacBook Pro 15-inch Touch":     {"MacBookPro13,3", "MacBookPro14,3", "MacBookPro15,1", "MacBookPro15,3", ""},
	"MacBook Pro 16-inch":           {"MacBookPro16,1"},
	"MacBook Pro 17-inch":           {"MacBookPro1,2", "MacBookPro2,1", "MacBookPro5,2", "MacBookPro6,1", "MacBookPro8,3"},
	"MacBook Pro 15/17-inch":        {"MacBookPro3,1", "MacBookPro4,1"},
	"Power Macintosh/Mac Server G3": {"PowerMac1,1"},
	"Power Macintosh/Mac Server G4": {"PowerMac1,2", "PowerMac3,1", "PowerMac3,3", "PowerMac3,4", "PowerMac3,5", "PowerMac3,6", "PowerMac5,1"},
	"Power Macintosh G5":            {"PowerMac7,2", "PowerMac7,3", "PowerMac9,1", "PowerMac11,2"},
	"PowerBook G3":                  {"PowerBook1,1", "PowerBook3,1"},
	"PowerBook G4":                  {"PowerBook3,2", "PowerBook3,3", "PowerBook3,4", "PowerBook3,5", "PowerBook5,1", "PowerBook5,2", "PowerBook5,3", "PowerBook5,4", "PowerBook5,5", "PowerBook5,6", "PowerBook5,7", "PowerBook5,8", "PowerBook5,9", "PowerBook6,1", "PowerBook6,2"},
	"Xserve G4":                     {"RackMac1,1"},
	"Xserve G5":                     {"RackMac3,1"},
	"Xserve Intel":                  {"Xserve1,1", "Xserve2,1", "Xserve3,1"},
}

var appleReverse map[string]string

func init() {
	appleReverse = make(map[string]string)

	for h, ids := range appleModels {
		for _, id := range ids {
			u := strings.ToUpper(id)

			existing, found := appleReverse[u]
			if found {
				fmt.Printf("Apple ID '%s' belongs to both '%s' and '%s'\n", id, h, existing)
			}

			appleReverse[u] = h
		}
	}
}

func appleHumanModel(id string) string {
	human, found := appleReverse[strings.ToUpper(id)]
	if found {
		return human
	}

	return id
}
