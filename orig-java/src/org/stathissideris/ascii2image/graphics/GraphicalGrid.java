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

	public int getCellMinX(TextGrid.Cell cell) {
		return getCellMinX(cell, getCellWidth());
	}

	public static int getCellMinX(TextGrid.Cell cell, int cellXSize) {
		return cell.x * cellXSize;
	}

	public int getCellMidX(TextGrid.Cell cell) {
		return getCellMidX(cell, getCellWidth());
	}

	public static int getCellMidX(TextGrid.Cell cell, int cellXSize) {
		return cell.x * cellXSize + cellXSize / 2;
	}

	public int getCellMaxX(TextGrid.Cell cell) {
		return getCellMaxX(cell, getCellWidth());
	}

	public static int getCellMaxX(TextGrid.Cell cell, int cellXSize) {
		return cell.x * cellXSize + cellXSize;
	}

	public int getCellMinY(TextGrid.Cell cell) {
		return getCellMinY(cell, getCellHeight());
	}

	public static int getCellMinY(TextGrid.Cell cell, int cellYSize) {
		return cell.y * cellYSize;
	}

	public int getCellMidY(TextGrid.Cell cell) {
		return getCellMidY(cell, getCellHeight());
	}

	public static int getCellMidY(TextGrid.Cell cell, int cellYSize) {
		return cell.y * cellYSize + cellYSize / 2;
	}

	public int getCellMaxY(TextGrid.Cell cell) {
		return getCellMaxY(cell, getCellHeight());
	}

	public static int getCellMaxY(TextGrid.Cell cell, int cellYSize) {
		return cell.y * cellYSize + cellYSize;
	}

	public TextGrid.Cell getCellFor(ShapePoint point) {
		if (point == null) {
			throw new IllegalArgumentException("ShapePoint cannot be null");
		}
		//TODO: the fake grid is a problem
		TextGrid g = new TextGrid();
		return g.new Cell((int) point.x / getCellWidth(), (int) point.y / getCellHeight());
	}

}