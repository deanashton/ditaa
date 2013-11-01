package org.stathissideris.ascii2image.graphics;

import org.stathissideris.ascii2image.core.RenderingOptions;
import org.stathissideris.ascii2image.text.TextGrid;

public class GraphicalGrid {
	private int width;
	private int height;
	private int cellWidth;
	private int cellHeight;

	public GraphicalGrid(TextGrid grid, RenderingOptions options) {
		cellWidth = options.getCellWidth();
		cellHeight = options.getCellHeight();

		width = grid.getWidth() * cellWidth;
		height = grid.getHeight() * cellHeight;
	}

	public int getWidth() {
		return width;
	}

	public int getHeight() {
		return height;
	}

	public int getCellWidth() {
		return cellWidth;
	}

	public int getCellHeight() {
		return cellHeight;
	}
}