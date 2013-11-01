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
package org.stathissideris.ascii2image.text;

import java.util.ArrayList;
import java.util.Iterator;

/**
 *
 * @author Efstathios Sideris
 */
public class GridPatternGroup extends ArrayList<GridPattern> {
	public boolean areAllMatchedBy(TextGrid grid) {
		Iterator<GridPattern> it = iterator();
		while (it.hasNext()) {
			GridPattern pattern = it.next();
			if (!pattern.isMatchedBy(grid)) {
				return false;
			}
		}
		return true;
	}

	public boolean isAnyMatchedBy(TextGrid grid) {
		Iterator<GridPattern> it = iterator();
		while (it.hasNext()) {
			GridPattern pattern = it.next();
			if (pattern.isMatchedBy(grid)) {
				return true;
			}
		}
		return false;
	}

	public void add(GridPattern... patterns) {
		for (GridPattern p : patterns) {
			add(p);
		}
	}

	//TODO: define criteria for on-line type?

	public static final GridPatternGroup cornerCriteria = new GridPatternGroup();
	public static final GridPatternGroup normalCornerCriteria = new GridPatternGroup();
	public static final GridPatternGroup roundCornerCriteria = new GridPatternGroup();

	public static final GridPatternGroup corner1Criteria = new GridPatternGroup();
	public static final GridPatternGroup corner2Criteria = new GridPatternGroup();
	public static final GridPatternGroup corner3Criteria = new GridPatternGroup();
	public static final GridPatternGroup corner4Criteria = new GridPatternGroup();

	public static final GridPatternGroup normalCorner1Criteria = new GridPatternGroup();
	public static final GridPatternGroup normalCorner2Criteria = new GridPatternGroup();
	public static final GridPatternGroup normalCorner3Criteria = new GridPatternGroup();
	public static final GridPatternGroup normalCorner4Criteria = new GridPatternGroup();

	public static final GridPatternGroup roundCorner1Criteria = new GridPatternGroup();
	public static final GridPatternGroup roundCorner2Criteria = new GridPatternGroup();
	public static final GridPatternGroup roundCorner3Criteria = new GridPatternGroup();
	public static final GridPatternGroup roundCorner4Criteria = new GridPatternGroup();

	public static final GridPatternGroup intersectionCriteria = new GridPatternGroup();
	public static final GridPatternGroup TCriteria = new GridPatternGroup();
	public static final GridPatternGroup inverseTCriteria = new GridPatternGroup();
	public static final GridPatternGroup KCriteria = new GridPatternGroup();
	public static final GridPatternGroup inverseKCriteria = new GridPatternGroup();

	public static final GridPatternGroup crossCriteria = new GridPatternGroup();

	public static final GridPatternGroup stubCriteria = new GridPatternGroup();
	public static final GridPatternGroup verticalLinesEndCriteria = new GridPatternGroup();
	public static final GridPatternGroup horizontalLinesEndCriteria = new GridPatternGroup();
	public static final GridPatternGroup linesEndCriteria = new GridPatternGroup();

	public static final GridPatternGroup crossOnLineCriteria = new GridPatternGroup();
	public static final GridPatternGroup horizontalCrossOnLineCriteria = new GridPatternGroup();
	public static final GridPatternGroup verticalCrossOnLineCriteria = new GridPatternGroup();

	public static final GridPatternGroup starOnLineCriteria = new GridPatternGroup();
	public static final GridPatternGroup horizontalStarOnLineCriteria = new GridPatternGroup();
	public static final GridPatternGroup verticalStarOnLineCriteria = new GridPatternGroup();

	public static final GridPatternGroup loneDiagonalCriteria = new GridPatternGroup();

	static {
		crossCriteria.add(new GridPattern(".6.", "4+8", ".2."));
		KCriteria.add(new GridPattern(".6.", "%4+8", ".2."));
		inverseKCriteria.add(new GridPattern(".6.", "4+%8", ".2."));
		TCriteria.add(new GridPattern(".%6.", "4+8", ".2."));
		inverseTCriteria.add(new GridPattern(".6.", "4+8", ".%2."));

		// ****** normal corners *******

		normalCorner1Criteria.add(new GridPattern(".[.", "~+(", ".^."));
		normalCorner2Criteria.add(new GridPattern(".[.", "(+~", ".^."));
		normalCorner3Criteria.add(new GridPattern(".^.", "(+~", ".[."));
		normalCorner4Criteria.add(new GridPattern(".^.", "~+(", ".[."));

		// ******* round corners *******

		roundCorner1Criteria.add(new GridPattern(".[.", "~/4", ".2."));
		roundCorner2Criteria.add(new GridPattern(".[.", "4\\~", ".2."));
		roundCorner3Criteria.add(new GridPattern(".6.", "4/~", ".[."));
		roundCorner4Criteria.add(new GridPattern(".6.", "~\\8", ".[."));

		//stubs

		stubCriteria.add(new GridPattern("!^!", "!+!", ".!."));
		stubCriteria.add(new GridPattern("!^!", "!+!", ".-."));
		stubCriteria.add(new GridPattern("!!.", "(+!", "!!."));
		stubCriteria.add(new GridPattern("!!.", "(+|", "!!."));
		stubCriteria.add(new GridPattern(".!.", "!+!", "!^!"));
		stubCriteria.add(new GridPattern(".-.", "!+!", "!^!"));
		stubCriteria.add(new GridPattern(".!!", "!+(", ".!!"));
		stubCriteria.add(new GridPattern(".!!", "|+(", ".!!"));

		// ****** ends of lines ******
		verticalLinesEndCriteria.add(new GridPattern(".^.", ".|.", ".!."));
		verticalLinesEndCriteria.add(new GridPattern(".^.", ".|.", ".-."));
		horizontalLinesEndCriteria.add(new GridPattern("...", "(-!", "..."));
		horizontalLinesEndCriteria.add(new GridPattern("...", "(-|", "..."));
		verticalLinesEndCriteria.add(new GridPattern(".!.", ".|.", ".^."));
		verticalLinesEndCriteria.add(new GridPattern(".-.", ".|.", ".^."));
		horizontalLinesEndCriteria.add(new GridPattern("...", "!-(", "..."));
		horizontalLinesEndCriteria.add(new GridPattern("...", "|-(", "..."));

		// ****** others *******

		horizontalCrossOnLineCriteria.add(new GridPattern("...", "(+(", "..."));
		verticalCrossOnLineCriteria.add(new GridPattern(".^.", ".+.", ".^."));
		horizontalStarOnLineCriteria.add(new GridPattern("...", "(*(", "..."));
		horizontalStarOnLineCriteria.add(new GridPattern("...", "!*(", "..."));
		horizontalStarOnLineCriteria.add(new GridPattern("...", "(*!", "..."));
		verticalStarOnLineCriteria.add(new GridPattern(".^.", ".*.", ".^."));
		verticalStarOnLineCriteria.add(new GridPattern(".!.", ".*.", ".^."));
		verticalStarOnLineCriteria.add(new GridPattern(".^.", ".*.", ".!."));
		loneDiagonalCriteria.add(new GridPattern(".%6%7", "%4/%8", "%3%2."));
		loneDiagonalCriteria.add(new GridPattern("%1%6.", "%4\\%8", ".%2%5"));

		//groups

		intersectionCriteria.addAll(crossCriteria);
		intersectionCriteria.addAll(KCriteria);
		intersectionCriteria.addAll(TCriteria);
		intersectionCriteria.addAll(inverseKCriteria);
		intersectionCriteria.addAll(inverseTCriteria);

		normalCornerCriteria.addAll(normalCorner1Criteria);
		normalCornerCriteria.addAll(normalCorner2Criteria);
		normalCornerCriteria.addAll(normalCorner3Criteria);
		normalCornerCriteria.addAll(normalCorner4Criteria);

		roundCornerCriteria.addAll(roundCorner1Criteria);
		roundCornerCriteria.addAll(roundCorner2Criteria);
		roundCornerCriteria.addAll(roundCorner3Criteria);
		roundCornerCriteria.addAll(roundCorner4Criteria);

		corner1Criteria.addAll(normalCorner1Criteria);
		corner1Criteria.addAll(roundCorner1Criteria);

		corner2Criteria.addAll(normalCorner2Criteria);
		corner2Criteria.addAll(roundCorner2Criteria);

		corner3Criteria.addAll(normalCorner3Criteria);
		corner3Criteria.addAll(roundCorner3Criteria);

		corner4Criteria.addAll(normalCorner4Criteria);
		corner4Criteria.addAll(roundCorner4Criteria);

		cornerCriteria.addAll(normalCornerCriteria);
		cornerCriteria.addAll(roundCornerCriteria);

		crossOnLineCriteria.addAll(horizontalCrossOnLineCriteria);
		crossOnLineCriteria.addAll(verticalCrossOnLineCriteria);

		starOnLineCriteria.addAll(horizontalStarOnLineCriteria);
		starOnLineCriteria.addAll(verticalStarOnLineCriteria);

		linesEndCriteria.addAll(horizontalLinesEndCriteria);
		linesEndCriteria.addAll(verticalLinesEndCriteria);
		linesEndCriteria.addAll(stubCriteria);

	}
}
