import { ColumnDef } from "@tanstack/react-table";
import * as _ from "lodash";

// Columns built with createColumnHelper have heterogeneous value types, so the
// value type parameter has to be `any` (as recommended by the TanStack docs).
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type TableColumns<T> = ColumnDef<T, any>[];

const PLACEHOLDER_ROWS = 3;

export default class TableUtils {
    // Replace every column accessor and cell renderer so that the table can be
    // rendered with empty placeholder rows while data is loading
    static GetPlaceholderColumns = function <T>(columns: TableColumns<T>): TableColumns<T> {
        return columns.map((col) => ({
            ...col,
            accessorFn: () => "",
            cell: () => <span className="skeleton inline-block h-3 w-3/4 align-middle"></span>,
        }));
    };

    static GetPlaceholderRows = function <T>(): T[] {
        return Array.from({ length: PLACEHOLDER_ROWS }, () => ({}) as T);
    };

    static getPaginationItems = function (
        currentPage: number,
        pageCount: number,
        maxLength: number
    ) {
        let res: Array<number> = [];

        if (pageCount <= maxLength) {
            return _.range(0, pageCount)
        } else {
            if (currentPage < 3) {
                res = _.range(3)
                res.push(NaN)

                res.push(..._.range(pageCount - 3, pageCount))
            } else if (currentPage >= pageCount - 3) {
                res.push(0)
                res.push(NaN)

                res.push(..._.range(pageCount - 3, pageCount))
            } else {
                res.push(0)
                res.push(NaN)

                res.push(..._.range(currentPage - 1, currentPage + 2))
                res.push(NaN)

                res.push(pageCount - 1)
            }

            return res
        }
    }
}
