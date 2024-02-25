import React from "react";
import {Column} from "react-table";
import * as _ from "lodash";

export default class TableUtils {
    static GetPlaceholderColumns = function<T extends Object>(columns :Array<Column<T>>) :Array<Column<T>> {
        let copy = _.cloneDeep(columns);
        return copy.map((col :Column<T>) => {
            col.accessor = () => "";
            //@ts-ignore
            col.cell = (cell :any) => <div className="px-2 py-3"><div className="animate-pulse rounded-full h-2 w-3/4 bg-slate-300 dark:bg-slate-600"></div></div>;
            return col;
        });
    };

    static getPaginationItems = function(
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
            } else if (currentPage >= pageCount - 3){
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
