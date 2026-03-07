package document

import (
	"Qingyu_backend/models/writer"
)

func rowToParagraph(row *writer.DocumentContent) writer.Paragraph {
	return writer.Paragraph{
		ID:          row.ID.Hex(),
		DocumentID:  row.DocumentID.Hex(),
		Order:       row.ParagraphOrder,
		Content:     row.Content,
		ContentType: row.ContentType,
		Version:     row.Version,
		UpdatedAt:   row.UpdatedAt,
	}
}

func paragraphToDTO(paragraph writer.Paragraph) ParagraphContent {
	return ParagraphContent{
		ParagraphID: paragraph.ID,
		Order:       paragraph.Order,
		Content:     paragraph.Content,
		ContentType: paragraph.ContentType,
		Version:     paragraph.Version,
		UpdatedAt:   paragraph.UpdatedAt,
	}
}

func dtoToParagraph(input ParagraphContent, defaultOrder int) writer.Paragraph {
	paragraph := writer.Paragraph{
		ID:          input.ParagraphID,
		Order:       input.Order,
		Content:     input.Content,
		ContentType: input.ContentType,
		Version:     input.Version,
		UpdatedAt:   input.UpdatedAt,
	}
	paragraph.Normalize(defaultOrder)
	return paragraph
}
