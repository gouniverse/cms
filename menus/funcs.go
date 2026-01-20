package cms

import (
	"context"
	"log"
	"slices"
	"sort"
	"strconv"

	"github.com/dracory/entitystore"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func getChildren(data []map[string]interface{}, parentID string) []map[string]interface{} {
	children := []map[string]interface{}{}
	//sequences := []string{}
	for _, node := range data {
		nodeParentID := ""
		//sequence := ""
		if keyExists(node, "parent_id") {
			nodeParentID = node["parent_id"].(string)
		}
		//if keyExists(node, "sequence") {
		//sequence = node["sequence"].(string)
		//}
		if nodeParentID == parentID {
			//sequences = append(sequences, sequence)
			children = append(children, node)
		}
	}

	sort.SliceStable(children, func(i, j int) bool {
		a, _ := strconv.ParseInt(children[i]["sequence"].(string), 10, 64)
		b, _ := strconv.ParseInt(children[j]["sequence"].(string), 10, 64)
		return a < b
	})

	return children
}

func buildTreeFromData(data []map[string]interface{}, parentID string) []map[string]interface{} {
	out := []map[string]interface{}{}

	roots := getChildren(data, parentID)

	for _, root := range roots {
		rootChildren := []map[string]interface{}{}
		rootID := root["id"].(string)
		children := getChildren(data, rootID)
		for childIndex, child := range children {
			childID := child["id"].(string)
			childrenTrees := buildTreeFromData(data, childID)
			rootChildren = append(rootChildren, childrenTrees...)
			children[childIndex]["children"] = rootChildren
		}
		root["children"] = children
		out = append(out, root)
	}

	return out
}

func (m UiManager) buildTree(menuID string) []map[string]interface{} {
	menuitems, err := m.entityStore.EntityListByAttribute(context.Background(), m.menuEntityType, "menu_id", menuID)

	if err != nil {
		log.Panicln("Menu items failed to be retrieved " + err.Error())
		return nil
	}

	nodeList := []map[string]interface{}{}
	for _, menuitem := range menuitems {
		itemID := menuitem.ID()
		itemName, _ := menuitem.GetString("name", "n/a")
		parentID, _ := menuitem.GetString("parent_id", "")
		sequence, _ := menuitem.GetString("sequence", "")
		target, _ := menuitem.GetString("target", "")
		url, _ := menuitem.GetString("url", "")
		pageID, _ := menuitem.GetString("page_id", "")
		item := map[string]interface{}{
			"id":        itemID,
			"parent_id": parentID,
			"sequence":  sequence,
			"name":      itemName,
			"page_id":   pageID,
			"url":       url,
			"target":    target,
		}
		nodeList = append(nodeList, item)
	}

	tree := buildTreeFromData(nodeList, "")

	return tree
}

func (m UiManager) pageMenusMenuItemsPagesDropdownList() (pagesDropdownList []map[string]string, errorMessage string) {
	pages, err := m.entityStore.EntityList(context.Background(), entitystore.EntityQueryOptions{
		EntityType: m.pageEntityType,
		Offset:     0,
		Limit:      200,
		SortBy:     "id",
		SortOrder:  "asc",
	})

	if err != nil {
		return pagesDropdownList, "Page list failed to be retrieved " + err.Error()
	}

	pagesDropdownList = make([]map[string]string, 0)

	mapPageIDTitle := map[string]string{}
	for _, page := range pages {
		title, err := page.GetString("title", "")

		if err != nil {
			return pagesDropdownList, "Page failed to be retrieved " + err.Error()
		}

		status, err := page.GetString("status", "")

		if err != nil {
			return pagesDropdownList, "Page failed to be retrieved " + err.Error()
		}

		mapPageIDTitle[page.ID()] = title + " (" + status + ")"
	}

	pageTitles := lo.Values(mapPageIDTitle)

	sort.Strings(pageTitles)

	pagesDropdownList = []map[string]string{}

	for _, title := range pageTitles {
		pageID, isFound := lo.FindKey(mapPageIDTitle, title)
		if !isFound {
			continue
		}
		pagesDropdownList = append(pagesDropdownList, map[string]string{
			"key":   pageID,
			"value": title,
		})
	}

	return pagesDropdownList, ""
}

// flattenTree flattens a JQTree data
func flattenTree(nodes []map[string]interface{}) []map[string]interface{} {
	flatTree := []map[string]interface{}{}
	for index, node := range nodes {
		children, hasChildren := node["children"]
		delete(node, "children")

		node["sequence"] = cast.ToString((index + 1))
		flatTree = append(flatTree, node)

		if !hasChildren {
			continue
		}

		childrenArray := children.([]interface{})
		childrenMapArray := []map[string]interface{}{}
		for _, child := range childrenArray {
			childMap := child.(map[string]interface{})
			childMap["parent_id"] = node["id"]
			childrenMapArray = append(childrenMapArray, childMap)
		}
		childNodesList := flattenTree(childrenMapArray)
		flatTree = append(flatTree, childNodesList...)
	}
	return flatTree
}

func (m UiManager) cleanMenuFromNonExistingMenuItems(menuID string, existingMenuItemIDs []string) (errorMessage string) {
	allMenuItems, err := m.entityStore.EntityListByAttribute(context.Background(), m.menuEntityType, "menu_id", menuID)

	if err != nil {
		return "Menu items failed to be fetched: " + err.Error()
	}

	// Delete old menu items
	for _, menuitem := range allMenuItems {
		exists := slices.Contains(existingMenuItemIDs, menuitem.ID())
		if !exists {
			m.entityStore.EntityDelete(context.Background(), menuitem.ID())
		}
	}

	return ""
}

func keyExists(decoded map[string]interface{}, key string) bool {
	val, ok := decoded[key]
	return ok && val != nil
}
