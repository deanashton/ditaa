/*
 * DiTAA - Diagrams Through Ascii Art
 *
 * Copyright (C) 2004 Efstathios Sideris
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA  02111-1307, USA.
 *
 */
package org.stathissideris.ascii2image.graphics;

import java.awt.Color;
import java.awt.Font;
import java.awt.geom.Rectangle2D;
import java.util.ArrayList;
import java.util.Iterator;

import org.stathissideris.ascii2image.core.ConversionOptions;
import org.stathissideris.ascii2image.core.Pair;
import org.stathissideris.ascii2image.text.AbstractionGrid;
import org.stathissideris.ascii2image.text.CellSet;
import org.stathissideris.ascii2image.text.TextGrid;
import org.stathissideris.ascii2image.text.TextGrid.CellStringPair;

/**
 *
 * @author Efstathios Sideris
 */
public class Diagram {

	private static final boolean DEBUG = true;
	private static final boolean DEBUG_VERBOSE = true;
	private static final boolean DEBUG_MAKE_SHAPES = false;

	private ArrayList<DiagramShape> shapes = new ArrayList<DiagramShape>();
	private ArrayList<CompositeDiagramShape> compositeShapes = new ArrayList<CompositeDiagramShape>();
	private ArrayList<DiagramText> textObjects = new ArrayList<DiagramText>();

	private GraphicalGrid ggrid = null;

	/**
	 *
	 * <p>
	 * An outline of the inner workings of this very important (and monstrous) constructor is presented here. Boundary processing is the first step of the
	 * process:
	 * </p>
	 *
	 * <ol>
	 * <li>Copy the grid into a work grid and remove all type-on-line and point markers from the work grid</li>
	 * <li>Split grid into distinct shapes by plotting the grid onto an AbstractionGrid and its getDistinctShapes() method.</li>
	 * <li>Find all the possible boundary sets of each of the distinct shapes. This can produce duplicate shapes (if the boundaries are the same when
	 * filling from the inside and the outside).</li>
	 * <li>Remove duplicate boundaries.</li>
	 * <li>Remove obsolete boundaries. Obsolete boundaries are the ones that are the sum of their parts when plotted as filled shapes. (see method
	 * removeObsoleteShapes())</li>
	 * <li>Seperate the found boundary sets to open, closed or mixed (See CellSet class on how its done).</li>
	 * <li>Are there any closed boundaries?
	 * <ul>
	 * <li>YES. Subtract all the closed boundaries from each of the open ones. That should convert the mixed shapes into open.</li>
	 * <li>NO. In this (harder) case, we use the method breakTrulyMixedBoundaries() of CellSet to break boundaries into open and closed shapes (would work
	 * in any case, but it's probably slower than the other method). This method is based on tracing from the lines' ends and splitting when we get to an
	 * intersection.</li>
	 * </ul>
	 * </li>
	 * <li>If we had to eliminate any mixed shapes, we seperate the found boundary sets again to open, closed or mixed.</li>
	 * </ol>
	 *
	 * <p>
	 * At this stage, the boundary processing is all complete and we proceed with using those boundaries to create the shapes:
	 * </p>
	 *
	 * <ol>
	 * <li>Create closed shapes.</li>
	 * <li>Create open shapes. That's when the line end corrections are also applied, concerning the positioning of the ends of lines see methods
	 * connectEndsToAnchors() and moveEndsToCellEdges() of DiagramShape.</li>
	 * <li>Assign color codes to closed shapes.</li>
	 * <li>Assing extended markup tags to closed shapes.</p>
	 * <li>Create arrowheads.</p>
	 * <li>Create point markers.</p>
	 * </ol>
	 *
	 * <p>
	 * Finally, the text processing occurs: [pending]
	 * </p>
	 *
	 */
	public Diagram(TextGrid textGrid, ConversionOptions options) {

		ggrid = new GraphicalGrid(textGrid, options.renderingOptions);

		TextGrid workGrid = new TextGrid(textGrid);
		workGrid.replaceTypeOnLine();
		workGrid.replacePointMarkersOnLine();
		if (DEBUG) {
			workGrid.printDebug();
		}

		int width = textGrid.getWidth();
		int height = textGrid.getHeight();

		//split distinct shapes using AbstractionGrid
		AbstractionGrid temp = new AbstractionGrid(workGrid, workGrid.getAllBoundaries());
		ArrayList<CellSet> boundarySetsStep1 = temp.getDistinctShapes();

		if (DEBUG) {
			System.out.println("******* Distinct shapes found using AbstractionGrid *******");
			for (CellSet set : boundarySetsStep1) {
				set.printAsGrid();
			}
			System.out.println("******* Same set of shapes after processing them by filling *******");
		}

		//Find all the boundaries by using the special version of the filling method
		//(fills in a different buffer than the buffer it reads from)
		ArrayList<CellSet> boundarySetsStep2 = new ArrayList<CellSet>();
		for (CellSet set : boundarySetsStep1) {
			//the fill buffer keeps track of which cells have been
			//filled already
			TextGrid fillBuffer = new TextGrid(width * 3, height * 3);

			for (int yi = 0; yi < height * 3; yi++) {
				for (int xi = 0; xi < width * 3; xi++) {
					if (fillBuffer.isBlank(xi, yi)) {

						TextGrid copyGrid = new AbstractionGrid(workGrid, set)
								.getCopyOfInternalBuffer();

						CellSet boundaries = copyGrid
								.findBoundariesExpandingFrom(copyGrid.new Cell(xi, yi));
						if (boundaries.size() == 0) {
							continue; //i'm not sure why these occur
						}
						boundarySetsStep2.add(boundaries.makeScaledOneThirdEquivalent());

						copyGrid = new AbstractionGrid(workGrid, set).getCopyOfInternalBuffer();
						CellSet filled = copyGrid.fillContinuousArea(copyGrid.new Cell(xi, yi),
								'*');
						fillBuffer.fillCellsWith(filled, '*');
						fillBuffer.fillCellsWith(boundaries, '-');

						if (DEBUG) {
							//System.out.println("Fill buffer:");
							//fillBuffer.printDebug();
							boundaries.makeScaledOneThirdEquivalent().printAsGrid();
							System.out.println("-----------------------------------");
						}

					}
				}
			}
		}

		if (DEBUG) {
			System.out.println("******* Removed duplicates *******");
		}

		boundarySetsStep2 = CellSet.removeDuplicateSets(boundarySetsStep2);

		if (DEBUG) {
			for (CellSet set : boundarySetsStep2) {
				set.printAsGrid();
			}
		}

		int originalSize = boundarySetsStep2.size();
		boundarySetsStep2 = CellSet.removeDuplicateSets(boundarySetsStep2);
		if (DEBUG) {
			System.out.println("******* Removed duplicates: there were " + originalSize
					+ " shapes and now there are " + boundarySetsStep2.size());
		}

		//split boundaries to open, closed and mixed

		if (DEBUG) {
			System.out.println("******* First evaluation of openess *******");
		}

		ArrayList<CellSet> open = new ArrayList<CellSet>();
		ArrayList<CellSet> closed = new ArrayList<CellSet>();
		ArrayList<CellSet> mixed = new ArrayList<CellSet>();

		for (CellSet set : boundarySetsStep2) {
			int type = set.getType(workGrid);
			if (type == CellSet.TYPE_CLOSED) {
				closed.add(set);
			} else if (type == CellSet.TYPE_OPEN) {
				open.add(set);
			} else if (type == CellSet.TYPE_MIXED) {
				mixed.add(set);
			}
			if (DEBUG) {
				if (type == CellSet.TYPE_CLOSED) {
					System.out.println("Closed boundaries:");
				} else if (type == CellSet.TYPE_OPEN) {
					System.out.println("Open boundaries:");
				} else if (type == CellSet.TYPE_MIXED) {
					System.out.println("Mixed boundaries:");
				}
				set.printAsGrid();
			}
		}

		boolean hadToEliminateMixed = false;

		if (mixed.size() > 0 && closed.size() > 0) {
			// mixed shapes can be eliminated by
			// subtracting all the closed shapes from them
			if (DEBUG) {
				System.out.println("******* Eliminating mixed shapes (basic algorithm) *******");
			}

			hadToEliminateMixed = true;

			//subtract from each of the mixed sets all the closed sets
			for (CellSet set : mixed) {
				for (CellSet closedSet : closed) {
					set.subtractSet(closedSet);
				}
				// this is necessary because some mixed sets produce
				// several distinct open sets after you subtract the
				// closed sets from them
				if (set.getType(workGrid) == CellSet.TYPE_OPEN) {
					boundarySetsStep2.remove(set);
					boundarySetsStep2.addAll(set.breakIntoDistinctBoundaries(workGrid));
				}
			}

		} else if (mixed.size() > 0 && closed.size() == 0) {
			// no closed shape exists, will have to
			// handle mixed shape on its own
			// an example of this case is the following:
			// +-----+
			// |  A  |C                 B
			// +  ---+-------------------
			// |     |
			// +-----+

			hadToEliminateMixed = true;

			if (DEBUG) {
				System.out.println("******* Eliminating mixed shapes (advanced algorithm for truly mixed shapes) *******");
			}

			for (CellSet set : mixed) {
				boundarySetsStep2.remove(set);
				boundarySetsStep2.addAll(set.breakTrulyMixedBoundaries(workGrid));
			}

		} else {
			if (DEBUG) {
				System.out.println("No mixed shapes found. Skipped mixed shape elimination step");
			}
		}

		if (hadToEliminateMixed) {
			if (DEBUG) {
				System.out.println("******* Second evaluation of openess *******");
			}

			//split boundaries again to open, closed and mixed
			open = new ArrayList<CellSet>();
			closed = new ArrayList<CellSet>();
			mixed = new ArrayList<CellSet>();

			for (CellSet set : boundarySetsStep2) {
				int type = set.getType(workGrid);
				if (type == CellSet.TYPE_CLOSED) {
					closed.add(set);
				} else if (type == CellSet.TYPE_OPEN) {
					open.add(set);
				} else if (type == CellSet.TYPE_MIXED) {
					mixed.add(set);
				}
				if (DEBUG) {
					if (type == CellSet.TYPE_CLOSED) {
						System.out.println("Closed boundaries:");
					} else if (type == CellSet.TYPE_OPEN) {
						System.out.println("Open boundaries:");
					} else if (type == CellSet.TYPE_MIXED) {
						System.out.println("Mixed boundaries:");
					}
					set.printAsGrid();
				}
			}
		}

		removeObsoleteShapes(workGrid, closed);

		boolean allCornersRound = false;
		if (options.processingOptions.areAllCornersRound()) {
			allCornersRound = true;
		}

		//make shapes from the boundary sets
		//make closed shapes
		if (DEBUG_MAKE_SHAPES) {
			System.out.println("***** MAKING SHAPES FROM BOUNDARY SETS *****");
			System.out.println("***** CLOSED: *****");
		}

		ArrayList<DiagramComponent> closedShapes = new ArrayList<DiagramComponent>();
		for (CellSet set : closed) {
			if (DEBUG_MAKE_SHAPES) {
				set.printAsGrid();
			}

			DiagramComponent shape = DiagramComponent.createClosedFromBoundaryCells(workGrid, set, ggrid
					.getCellWidth(), ggrid.getCellHeight(), allCornersRound);
			if (shape != null) {
				if (shape instanceof DiagramShape) {
					addToShapes((DiagramShape) shape);
					closedShapes.add(shape);
				} else if (shape instanceof CompositeDiagramShape) {
					addToCompositeShapes((CompositeDiagramShape) shape);
				}
			}
		}

		if (options.processingOptions.performSeparationOfCommonEdges()) {
			separateCommonEdges(closedShapes);
		}

		//make open shapes
		for (CellSet set : open) {
			if (set.size() == 1) { //single cell "shape"
				TextGrid.Cell cell = set.getFirst();
				if (!textGrid.cellContainsDashedLineChar(cell)) {
					DiagramShape shape = DiagramShape.createSmallLine(workGrid, cell, ggrid
							.getCellWidth(), ggrid.getCellHeight());
					if (shape != null) {
						addToShapes(shape);
						shape.connectEndsToAnchors(workGrid, ggrid);
					}
				}
			} else { //normal shape
				if (DEBUG) {
					System.out.println(set.getCellsAsString());
				}

				DiagramComponent shape = CompositeDiagramShape.createOpenFromBoundaryCells(workGrid,
						set, ggrid.getCellWidth(), ggrid.getCellHeight(), allCornersRound);

				if (shape != null) {
					if (shape instanceof CompositeDiagramShape) {
						addToCompositeShapes((CompositeDiagramShape) shape);
						((CompositeDiagramShape) shape).connectEndsToAnchors(workGrid, ggrid);
					} else if (shape instanceof DiagramShape) {
						addToShapes((DiagramShape) shape);
						((DiagramShape) shape).connectEndsToAnchors(workGrid, ggrid);
						((DiagramShape) shape).moveEndsToCellEdges(textGrid, ggrid);
					}
				}

			}
		}

		//assign color codes to shapes
		//TODO: text on line should not change its color

		for (TextGrid.CellColorPair pair : textGrid.findColorCodes()) {
			ShapePoint point = new ShapePoint(ggrid.getCellMidX(pair.cell), ggrid.getCellMidY(pair.cell));
			DiagramShape containingShape = findSmallestShapeContaining(point);

			if (containingShape != null) {
				containingShape.setFillColor(pair.color);
			}
		}

		//assign markup to shapes
		for (TextGrid.CellTagPair pair : textGrid.findMarkupTags()) {
			ShapePoint point = new ShapePoint(ggrid.getCellMidX(pair.cell), ggrid.getCellMidY(pair.cell));

			DiagramShape containingShape = findSmallestShapeContaining(point);

			//this tag is not within a shape, skip
			if (containingShape == null) {
				continue;
			}

			if (pair.tag.equals("d")) {
				loadCustomShape(containingShape, options, "d", DiagramShape.TYPE_DOCUMENT);
			} else if (pair.tag.equals("s")) {
				loadCustomShape(containingShape, options, "s", DiagramShape.TYPE_STORAGE);
			} else if (pair.tag.equals("io")) {
				loadCustomShape(containingShape, options, "io", DiagramShape.TYPE_IO);
			} else if (pair.tag.equals("c")) {
				loadCustomShape(containingShape, options, "c", DiagramShape.TYPE_DECISION);
			} else if (pair.tag.equals("mo")) {
				loadCustomShape(containingShape, options, "mo", DiagramShape.TYPE_MANUAL_OPERATION);
			} else if (pair.tag.equals("tr")) {
				loadCustomShape(containingShape, options, "tr", DiagramShape.TYPE_TRAPEZOID);
			} else if (pair.tag.equals("o")) {
				loadCustomShape(containingShape, options, "o", DiagramShape.TYPE_ELLIPSE);
			} else {
				CustomShapeDefinition def = options.processingOptions.getFromCustomShapes(pair.tag);
				containingShape.setType(DiagramShape.TYPE_CUSTOM);
				containingShape.setDefinition(def);
			}
		}

		//make arrowheads
		for (TextGrid.Cell cell : workGrid.findArrowheads()) {
			DiagramShape arrowhead = DiagramShape.createArrowhead(workGrid, cell, ggrid.getCellWidth(),
					ggrid.getCellHeight());
			if (arrowhead != null) {
				addToShapes(arrowhead);
			} else {
				System.err.println("Could not create arrowhead shape. Unexpected error.");
			}
		}

		//make point markers
		for (TextGrid.Cell cell : textGrid.getPointMarkersOnLine()) {
			DiagramShape mark = new DiagramShape();
			mark.addToPoints(new ShapePoint(ggrid.getCellMidX(cell), ggrid.getCellMidY(cell)));
			mark.setType(DiagramShape.TYPE_POINT_MARKER);
			mark.setFillColor(Color.white);
			shapes.add(mark);
		}

		removeDuplicateShapes();

		if (DEBUG) {
			System.out.println("Shape count: " + shapes.size());
		}
		if (DEBUG) {
			System.out.println("Composite shape count: " + compositeShapes.size());
		}

		//copy again
		workGrid = new TextGrid(textGrid);
		workGrid.removeNonText();

		// ****** handle text *******
		//break up text into groups
		TextGrid textGroupGrid = new TextGrid(workGrid);
		CellSet gaps = textGroupGrid.getAllBlanksBetweenCharacters();
		//kludge
		textGroupGrid.fillCellsWith(gaps, '|');
		CellSet nonBlank = textGroupGrid.getAllNonBlank();
		ArrayList<CellSet> textGroups = nonBlank.breakIntoDistinctBoundaries();
		if (DEBUG) {
			System.out.println(textGroups.size() + " text groups found");
		}

		Font font = FontMeasurer.instance().getFontFor(ggrid.getCellHeight());

		for (CellSet textGroupCellSet : textGroups) {
			TextGrid isolationGrid = new TextGrid(width, height);
			workGrid.copyCellsTo(textGroupCellSet, isolationGrid);

			ArrayList<CellStringPair> strings = isolationGrid.findStrings();
			for (TextGrid.CellStringPair pair : strings) {
				TextGrid.Cell cell = pair.cell;
				String string = pair.string;
				if (DEBUG) {
					System.out.println("Found string " + string);
				}
				TextGrid.Cell lastCell = isolationGrid.new Cell(cell.x + string.length() - 1, cell.y);

				int minX = ggrid.getCellMinX(cell);
				int y = ggrid.getCellMaxY(cell);
				int maxX = ggrid.getCellMaxX(lastCell);

				DiagramText textObject;
				if (FontMeasurer.instance().getWidthFor(string, font) > maxX - minX) { //does not fit horizontally
					Font lessWideFont = FontMeasurer.instance().getFontFor(maxX - minX, string);
					textObject = new DiagramText(minX, y, string, lessWideFont);
				} else {
					textObject = new DiagramText(minX, y, string, font);
				}

				textObject.centerVerticallyBetween(ggrid.getCellMinY(cell), ggrid.getCellMaxY(cell));

				//TODO: if the strings start with bullets they should be aligned to the left

				//position text correctly
				int otherStart = isolationGrid.otherStringsStartInTheSameColumn(cell);
				int otherEnd = isolationGrid.otherStringsEndInTheSameColumn(lastCell);
				if (0 == otherStart && 0 == otherEnd) {
					textObject.centerHorizontallyBetween(minX, maxX);
				} else if (otherEnd > 0 && otherStart == 0) {
					textObject.alignRightEdgeTo(maxX);
				} else if (otherEnd > 0 && otherStart > 0) {
					if (otherEnd > otherStart) {
						textObject.alignRightEdgeTo(maxX);
					} else if (otherEnd == otherStart) {
						textObject.centerHorizontallyBetween(minX, maxX);
					}
				}

				addToTextObjects(textObject);
			}
		}

		if (DEBUG) {
			System.out.println("Positioned text");
		}

		//correct the color of the text objects according
		//to the underlying color
		for (DiagramText textObject : getTextObjects()) {
			DiagramShape shape = findSmallestShapeIntersecting(textObject.getBounds());
			if (shape != null && shape.getFillColor() != null
					&& BitmapRenderer.isColorDark(shape.getFillColor())) {
				textObject.setColor(Color.white);
			}
		}

		//set outline to true for test within custom shapes
		for (DiagramShape shape : getAllDiagramShapes()) {
			if (shape.getType() == DiagramShape.TYPE_CUSTOM) {
				for (DiagramText textObject : getTextObjects()) {
					textObject.setHasOutline(true);
					textObject.setColor(DiagramText.DEFAULT_COLOR);
				}
			}
		}

		if (DEBUG) {
			System.out.println("Corrected color of text according to underlying color");
		}

	}

	private void loadCustomShape(DiagramShape containingShape, ConversionOptions options, String code,
			int fallbackType) {
		CustomShapeDefinition def = options.processingOptions.getFromCustomShapes(code);
		if (def == null) {
			containingShape.setType(fallbackType);
		} else {
			containingShape.setType(DiagramShape.TYPE_CUSTOM);
			containingShape.setDefinition(def);
		}
	}

	/**
	 * Returns a list of all DiagramShapes in the Diagram, including the ones within CompositeDiagramShapes
	 *
	 * @return
	 */
	public ArrayList<DiagramShape> getAllDiagramShapes() {
		ArrayList<DiagramShape> shapes = new ArrayList<DiagramShape>();
		shapes.addAll(getShapes());

		for (CompositeDiagramShape compShape : getCompositeShapes()) {
			shapes.addAll(compShape.getShapes());
		}
		return shapes;
	}

	/**
	 * Removes the sets from <code>sets</code>that are the sum of their parts when plotted as filled shapes.
	 *
	 * @return true if it removed any obsolete.
	 *
	 */
	private boolean removeObsoleteShapes(TextGrid grid, ArrayList<CellSet> sets) {
		if (DEBUG) {
			System.out.println("******* Removing obsolete shapes *******");
		}

		boolean removedAny = false;

		ArrayList<CellSet> filledSets = new ArrayList<CellSet>();

		Iterator<?> it;

		if (DEBUG_VERBOSE) {
			System.out.println("******* Sets before *******");
			for (CellSet set : sets) {
				set.printAsGrid();
			}
		}

		//make filled versions of all the boundary sets
		for (CellSet set : sets) {
			set = set.getFilledEquivalent(grid);
			if (set == null) {
				return false;
			} else {
				filledSets.add(set);
			}
		}

		ArrayList<Integer> toBeRemovedIndices = new ArrayList<Integer>();
		for (CellSet set : filledSets) {
			if (DEBUG_VERBOSE) {
				System.out.println("*** Deciding if the following should be removed:");
				set.printAsGrid();
			}

			//find the other sets that have common cells with set
			ArrayList<CellSet> common = new ArrayList<CellSet>();
			common.add(set);
			for (CellSet set2 : filledSets) {
				if (set != set2 && set.hasCommonCells(set2)) {
					common.add(set2);
				}
			}
			//it only makes sense for more than 2 sets
			if (common.size() == 2) {
				continue;
			}

			//find largest set
			CellSet largest = set;
			for (CellSet set2 : common) {
				if (set2.size() > largest.size()) {
					largest = set2;
				}
			}

			if (DEBUG_VERBOSE) {
				System.out.println("Largest:");
				largest.printAsGrid();
			}

			//see if largest is sum of others
			common.remove(largest);

			//make the sum set of the small sets on a grid
			TextGrid gridOfSmalls = new TextGrid(largest.getMaxX() + 2, largest.getMaxY() + 2);
			new CellSet();
			for (CellSet set2 : common) {
				if (DEBUG_VERBOSE) {
					System.out.println("One of smalls:");
					set2.printAsGrid();
				}
				gridOfSmalls.fillCellsWith(set2, '*');
			}
			if (DEBUG_VERBOSE) {
				System.out.println("Sum of smalls:");
				gridOfSmalls.printDebug();
			}
			TextGrid gridLargest = new TextGrid(largest.getMaxX() + 2, largest.getMaxY() + 2);
			gridLargest.fillCellsWith(largest, '*');

			int index = filledSets.indexOf(largest);
			if (gridLargest.equals(gridOfSmalls) && !toBeRemovedIndices.contains(new Integer(index))) {
				toBeRemovedIndices.add(new Integer(index));
				if (DEBUG) {
					System.out.println("Decided to remove set:");
					largest.printAsGrid();
				}
			} /*else if (DEBUG){
				System.out.println("This set WILL NOT be removed:");
				largest.printAsGrid();
				}*/
			//if(gridLargest.equals(gridOfSmalls)) toBeRemovedIndices.add(new Integer(index));
		}

		ArrayList<CellSet> setsToBeRemoved = new ArrayList<CellSet>();
		for (Integer i : toBeRemovedIndices) {
			setsToBeRemoved.add(sets.get(i));
		}

		for (CellSet set : setsToBeRemoved) {
			removedAny = true;
			sets.remove(set);
		}

		if (DEBUG_VERBOSE) {
			System.out.println("******* Sets after *******");
			for (CellSet set : sets) {
				set.printAsGrid();
			}
		}

		return removedAny;
	}

	private void separateCommonEdges(ArrayList<DiagramComponent> shapes) {

		float offset = ggrid.getMinimumOfCellDimension() / 5;

		ArrayList<ShapeEdge> edges = new ArrayList<ShapeEdge>();

		//get all adges
		for (DiagramComponent component : shapes) {
			DiagramShape shape = (DiagramShape) component;
			edges.addAll(shape.getEdges());
		}

		//group edges into pairs of touching edges
		ArrayList<Pair<ShapeEdge, ShapeEdge>> listOfPairs = new ArrayList<Pair<ShapeEdge, ShapeEdge>>();

		//all-against-all touching test for the edges
		int startIndex = 1; //skip some to avoid duplicate comparisons and self-to-self comparisons

		for (ShapeEdge edge1 : edges) {
			for (int k = startIndex; k < edges.size(); k++) {
				ShapeEdge edge2 = edges.get(k);

				if (edge1.touchesWith(edge2)) {
					listOfPairs.add(new Pair<ShapeEdge, ShapeEdge>(edge1, edge2));
				}
			}
			startIndex++;
		}

		ArrayList<ShapeEdge> movedEdges = new ArrayList<ShapeEdge>();

		//move equivalent edges inwards
		for (Pair<ShapeEdge, ShapeEdge> pair : listOfPairs) {
			if (!movedEdges.contains(pair.first)) {
				pair.first.moveInwardsBy(offset);
				movedEdges.add(pair.first);
			}
			if (!movedEdges.contains(pair.second)) {
				pair.second.moveInwardsBy(offset);
				movedEdges.add(pair.second);
			}
		}

	}

	//TODO: removes more than it should
	private void removeDuplicateShapes() {
		ArrayList<DiagramShape> originalShapes = new ArrayList<DiagramShape>();

		for (DiagramShape shape : shapes) {
			boolean isOriginal = true;
			for (DiagramShape originalShape : originalShapes) {
				if (shape.equals(originalShape)) {
					isOriginal = false;
				}
			}
			if (isOriginal) {
				originalShapes.add(shape);
			}
		}

		shapes.clear();
		shapes.addAll(originalShapes);
	}

	private DiagramShape findSmallestShapeContaining(ShapePoint point) {
		DiagramShape containingShape = null;
		for (DiagramShape shape : shapes) {
			if (shape.contains(point)) {
				if (containingShape == null) {
					containingShape = shape;
				} else {
					if (shape.isSmallerThan(containingShape)) {
						containingShape = shape;
					}
				}
			}
		}
		return containingShape;
	}

	private DiagramShape findSmallestShapeIntersecting(Rectangle2D rect) {
		DiagramShape intersectingShape = null;
		for (DiagramShape shape : shapes) {
			if (shape.intersects(rect)) {
				if (intersectingShape == null) {
					intersectingShape = shape;
				} else {
					if (shape.isSmallerThan(intersectingShape)) {
						intersectingShape = shape;
					}
				}
			}
		}
		return intersectingShape;
	}

	private void addToTextObjects(DiagramText shape) {
		textObjects.add(shape);
	}

	private void addToCompositeShapes(CompositeDiagramShape shape) {
		compositeShapes.add(shape);
	}

	private void addToShapes(DiagramShape shape) {
		shapes.add(shape);
	}

	public Iterator<DiagramShape> getShapesIterator() {
		return shapes.iterator();
	}

	public GraphicalGrid getGraphicalGrid() {
		return ggrid;
	}

	/**
	 * @return
	 */
	public int getHeight() {
		return ggrid.getHeight();
	}

	/**
	 * @return
	 */
	public int getWidth() {
		return ggrid.getWidth();
	}

	/**
	 * @return
	 */
	public int getCellWidth() {
		return ggrid.getCellWidth();
	}

	/**
	 * @return
	 */
	public int getCellHeight() {
		return ggrid.getCellHeight();
	}

	/**
	 * @return
	 */
	public ArrayList<CompositeDiagramShape> getCompositeShapes() {
		return compositeShapes;
	}

	/**
	 * @return
	 */
	public ArrayList<DiagramShape> getShapes() {
		return shapes;
	}

	/**
	 * @return
	 */
	public ArrayList<DiagramText> getTextObjects() {
		return textObjects;
	}

}
