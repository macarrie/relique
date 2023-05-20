import React, {useState} from "react";
import {useGlobalFilter, useSortBy, useTable} from "react-table";
import * as _ from "lodash";

function Table({title, filtered, sorted, refreshFunc, columns, data, actions} :any) {
    let [filter, setTableFilter] = useState("");

    const tableInstance = useTable(
        {
            columns: columns,
            data: data,
            // @ts-ignore
            disableSortBy: !sorted,
        },
        useGlobalFilter,
        useSortBy,
    );
    const {
        getTableProps,
        getTableBodyProps,
        headerGroups,
        rows,
        prepareRow,
        // @ts-ignore
        setGlobalFilter,
    } = tableInstance

    function updateFilter(value :string) {
        setTableFilter(value);
        _.debounce((value :any) => setGlobalFilter(value), 200)(value);
    }

    return (
        <div>
            <div className="flex flex-row px-4 py-3 items-center">
                <span className="flex-grow text-xl">{title}</span>
                {filtered === true && (
                    <div className="flex flex-row items-center rounded bg-slate-100 dark:bg-slate-800">
                        <input type="text" className="bg-transparent border-none text-sm focus:ring-0" value={filter} onChange={e => updateFilter(e.target.value)} placeholder="Quick search" />
                        {filter !== "" && (
                            <button type="button" className="button button-small button-text" onClick={() => updateFilter("")}>
                                <i className="ri-close-line"></i>
                            </button>
                        )}
                    </div>
                )}
                {refreshFunc && (
                    <button type={"button"} className="button button-small button-text ml-2" onClick={refreshFunc}>
                        <div className="hover:animate-spin">
                            <i className="text-lg ri-refresh-line"></i>
                        </div>
                    </button>
                )}
                {actions}
            </div>
            <table className="table table-auto w-full" {...getTableProps()}>
                <thead>
                {
                    headerGroups.map((headerGroup :any) => (
                        <tr {...headerGroup.getHeaderGroupProps()}>
                            {
                                headerGroup.headers.map((column :any) => (
                                    <th {...column.getHeaderProps(column.getSortByToggleProps())}>
                                        <div className="flex flex-row items-center">
                                            { column.render('Header') }
                                            <span className="text-slate-400 dark:text-slate-500">
                                                {column.isSorted ? (column.isSortedDesc ? <i className="ri-sort-desc"></i> : <i className="ri-sort-asc"></i>) : ''}
                                            </span>
                                        </div>
                                    </th>
                                ))}
                        </tr>
                    ))}
                </thead>
                {data.length === 0 ? (
                    <tbody className="text-3xl text-center italic text-slate-300 dark:text-slate-700">
                        <tr><td className="py-8" colSpan={100}>No elements</td></tr>
                    </tbody>
                ) : (
                    <tbody {...getTableBodyProps()}>
                    {
                        rows.map((row :any) => {
                            prepareRow(row)
                            return (
                                <tr {...row.getRowProps()}>
                                    {
                                        row.cells.map((cell :any) => {
                                            return (
                                                <td {...cell.getCellProps()}>
                                                    {
                                                        cell.render('Cell')}
                                                </td>
                                            )
                                        })}
                                </tr>
                            )
                        })
                    }
                    </tbody>
                )}
            </table>
        </div>
    );
}
export default Table;
