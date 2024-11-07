package cms

import (
	"strings"
)

// RenderBlocks renders the blocks in a string
func (cms *Cms) ContentRenderBlocks(content string) (string, error) {
	blockIDs := ContentFindIdsByPatternPrefix(content, "BLOCK")

	var err error
	for _, blockID := range blockIDs {
		content, err = cms.ContentRenderBlockByID(content, blockID)

		if err != nil {
			return content, err
		}
	}

	return content, nil
}

// ContentRenderBlockByID renders the block specified by the ID in a content
// if the blockID is empty or not found the initial content is returned
func (cms *Cms) ContentRenderBlockByID(content string, blockID string) (string, error) {
	if blockID == "" {
		return content, nil
	}

	blockContent, err := cms.findBlockContent(blockID)

	if err != nil {
		return "", err
	}

	content = strings.ReplaceAll(content, "[[BLOCK_"+blockID+"]]", blockContent)
	content = strings.ReplaceAll(content, "[[ BLOCK_"+blockID+" ]]", blockContent)

	return content, nil
}

func (cms *Cms) findBlockContent(blockID string) (string, error) {
	block, err := cms.BlockFindByID(blockID)

	if err != nil {
		return "", err
	}

	var blockContent string

	if block == nil {
		blockContent = ""
	} else {
		blockStatus := block.Status()

		if blockStatus == "active" {
			blockContent = block.Content()
		}
	}

	return blockContent, nil
}
