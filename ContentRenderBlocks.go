package cms

import (
	"log"
	"strings"
)

// RenderBlocks renders the blocks in a string
func (cms *Cms) ContentRenderBlocks(content string) string {
	blockIDs := ContentFindIdsByPatternPrefix(content, "BLOCK")

	for _, blockID := range blockIDs {
		content = cms.ContentRenderBlockByID(content, blockID)
	}

	return content
}

// RenderBlockByID renders the block specified by the ID in a content
// if the blockID is empty or not found the initial content is returned
func (cms *Cms) ContentRenderBlockByID(content string, blockID string) string {
	if blockID == "" {
		return content
	}

	blockContent := cms.findBlockContent(blockID)
	content = strings.ReplaceAll(content, "[[BLOCK_"+blockID+"]]", blockContent)
	content = strings.ReplaceAll(content, "[[ BLOCK_"+blockID+" ]]", blockContent)
	return content
}

func (cms *Cms) findBlockContent(blockID string) string {
	block, _ := cms.EntityStore.EntityFindByID(blockID)

	var blockContent string

	if block == nil {
		log.Println("Block " + blockID + " not found")
		blockContent = ""
	} else {
		blockStatus, _ := block.GetString("status", "")
		if blockStatus == "active" {
			blockContent, _ = block.GetString("content", "")
		}
	}

	return blockContent
}
