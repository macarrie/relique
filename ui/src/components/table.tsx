import React, {useEffect, useState} from "react";
import {
    flexRender,
    SortingState,
    getCoreRowModel,
    getSortedRowModel,
    getFilteredRowModel,
    getPaginationRowModel,
    useReactTable,
} from "@tanstack/react-table";
import * as _ from "lodash";
import TableUtils from "../utils/table";
import {range} from "lodash";
import PaginationButton from "./pagination_button";

function Table({
                   title,
                   filtered,
                   sorted,
                   paginated,
                   columns,
                   data,
                   loading,
                   defaultPageSize,
                   refreshFunc,
                   actions,
} :any) {
    const placeholderColumns = TableUtils.GetPlaceholderColumns(columns)

    // For display (without debounce)
    let [filter, setDisplayFilter] = useState("");
    // For table filtering (updated with debounce)
    let [globalFilter, setGlobalFilter] = useState("");
    let [sorting, setSorting] = React.useState<SortingState>([])

    const table = useReactTable(
        {
            columns: loading ? placeholderColumns : columns,
            data: loading ? [{}, {}, {}] : data,
            state: {
                globalFilter,
                sorting,
            },
            onGlobalFilterChange: updateFilter,
            onSortingChange: setSorting,
            getCoreRowModel: getCoreRowModel(),
            getSortedRowModel: getSortedRowModel(),
            getFilteredRowModel: getFilteredRowModel(),
            getPaginationRowModel: getPaginationRowModel(),
            enableGlobalFilter: filtered,
            enableSorting: sorted,
        },
    );

    useEffect(() => {
        table.setPageSize(defaultPageSize)
    }, [table, defaultPageSize])

    function updateFilter(value :string) {
        setDisplayFilter(value)
        _.debounce((value :any) => setGlobalFilter(value), 200)(value);
    }

    return (
        <div>
            <div className="flex flex-row px-4 py-3 items-center">
                <span className="flex-grow text-xl">{title}</span>
                {filtered === true && (
                    <div className="flex flex-row items-center rounded bg-slate-100 dark:bg-slate-800">
                        <input type="text" className="bg-transparent border-none text-sm focus:ring-0" value={filter}
                               onChange={e => updateFilter(e.target.value)} placeholder="Quick search"/>
                        {filter !== "" && (
                            <button type="button" className="button button-small button-text"
                                    onClick={() => updateFilter("")}>
                                <i className="ri-close-line"></i>
                            </button>
                        )}
                    </div>
                )}
                <button type={"button"} className="button button-small button-text ml-2"
                        onClick={() => refreshFunc()}>
                    <div className="hover:animate-spin">
                        <i className="text-lg ri-refresh-line"></i>
                    </div>
                </button>
                {actions}
            </div>
            <table className="table table-auto w-full">
                <thead>
                {table.getHeaderGroups().map(headerGroup => (
                    <tr key={headerGroup.id}>
                        {headerGroup.headers.map(header => (
                            <th key={header.id}>
                                {header.isPlaceholder ? null : (
                                    <div
                                        {...{
                                            className: header.column.getCanSort()
                                                ? 'cursor-pointer select-none'
                                                : '',
                                            onClick: header.column.getToggleSortingHandler(),
                                        }}
                                    >
                                        <div className="flex flex-row items-center">
                                            {flexRender(
                                                header.column.columnDef.header,
                                                header.getContext()
                                            )}

                                            <span className="text-slate-400 dark:text-slate-500">
                                            {{
                                                asc: <i className="ri-sort-asc"></i>,
                                                desc: <i className="ri-sort-desc"></i>,
                                            }[header.column.getIsSorted() as string] ?? null}
                                            </span>
                                        </div>
                                    </div>
                                )}
                            </th>
                        ))}
                    </tr>
                ))}
                </thead>
                {table.getTotalSize() === 0 ? (
                    <tbody className="text-3xl text-center italic text-slate-300 dark:text-slate-700">
                    <tr>
                        <td className="py-8" colSpan={100}>No elements</td>
                    </tr>
                    </tbody>
                ) : (
                    <tbody>
                    {table.getRowModel().rows.map(row => (
                        <tr key={row.id}>
                            {row.getVisibleCells().map(cell => (
                                <td key={cell.id}>
                                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                                </td>
                            ))}
                        </tr>
                    ))}
                    </tbody>
                )}
                <tfoot>
                {table.getFooterGroups().map(footerGroup => (
                    <tr key={footerGroup.id}>
                        {footerGroup.headers.map(header => (
                            <th key={header.id}>
                                {header.isPlaceholder
                                    ? null
                                    : flexRender(
                                        header.column.columnDef.footer,
                                        header.getContext()
                                    )}
                            </th>
                        ))}
                    </tr>
                ))}
                </tfoot>
            </table>

            {(paginated && !loading) && (
                <div className="flex flex-row px-4 pt-2">
                <span className="flex items-center gap-2">
                    <b>
                        {table.getPrePaginationRowModel().rows.length}
                    </b>
                    <div>
                        elements displayed
                    </div>
                </span>
                    <div className="flex flex-grow flex-row justify-center items-center">
                        <div className="space-x-1">
                            <PaginationButton
                                    onClick={() => table.setPageIndex(0)} disabled={!table.getCanPreviousPage()}>
                                <i className="ri-arrow-left-double-line"></i>
                            </PaginationButton>
                            <PaginationButton
                                    onClick={() => table.previousPage()} disabled={!table.getCanPreviousPage()}>
                                <i className="ri-arrow-left-s-line"></i>
                            </PaginationButton>
                            {range(0, table.getPageCount()).map((page: number) => (
                                <PaginationButton
                                    key={page}
                                    active={table.getState().pagination.pageIndex === page}
                                    onClick={() => table.setPageIndex(page)}>
                                    {page + 1}
                                </PaginationButton>
                            ))}
                            <PaginationButton
                                    onClick={() => table.nextPage()} disabled={!table.getCanNextPage()}>
                                <i className="ri-arrow-right-s-line"></i>
                            </PaginationButton>
                            <PaginationButton onClick={() => table.setPageIndex(table.getPageCount() - 1)}
                                              disabled={!table.getCanNextPage()}>
                                <i className="ri-arrow-right-double-line"></i>
                            </PaginationButton>
                        </div>
                    </div>
                    <div>
                        <select
                            value={table.getState().pagination.pageSize}
                            onChange={e => {
                                table.setPageSize(Number(e.target.value))
                            }}
                            className="form-select bg-transparent text-sm border-none">
                            {[10, 25, 50, 100].map(size => (
                                <option key={size} value={size}>
                                    {size} per page
                                </option>
                            ))}
                        </select>
                    </div>
                </div>
            )}
        </div>
    );
}
export default Table;
