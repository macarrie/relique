import React, { useState } from "react";
import {
    flexRender,
    SortingState,
    getCoreRowModel,
    getSortedRowModel,
    getPaginationRowModel,
    PaginationState,
    useReactTable,
    getFilteredRowModel,
    FilterFn,
} from "@tanstack/react-table";
// A TanStack fork of Kent C. Dodds' match-sorter library that provides ranking information
import {
    RankingInfo,
    rankItem,
} from '@tanstack/match-sorter-utils'

import * as _ from "lodash";
import TableUtils from "../utils/table";
import PaginationButton from "./pagination_button";
import Const from "../types/const";
import DebouncedInput from "./debounced_input";

declare module '@tanstack/react-table' {
    //add fuzzy filter to the filterFns
    interface FilterFns {
        fuzzy: FilterFn<unknown>
    }
    interface FilterMeta {
        itemRank: RankingInfo
    }
}

// Define a custom fuzzy filter function that will apply ranking info to rows (using match-sorter utils)
const fuzzyFilter: FilterFn<any> = (row, columnId, value, addMeta) => {
    // Rank the item
    const itemRank = rankItem(row.getValue(columnId), value)

    // Store the itemRank info
    addMeta({
        itemRank,
    })

    // Return if the item should be filtered in/out
    return itemRank.passed
}

function Table({
    title = "",
    actions = true,
    custom_actions = [] as React.ReactNode,
    filtered = true,
    paginated = true,
    sorted = true,
    columns,
    data,
    defaultPageSize = Const.DEFAULT_PAGE_SIZE,
}: any) {
    const [pagination, setPagination] = useState<PaginationState>({
        pageIndex: 0,
        pageSize: defaultPageSize,
    })
    const [sorting, setSorting] = useState<SortingState>([]);
    const [globalFilter, setGlobalFilter] = useState('');

    let paginationOptions = paginated ? {
        getPaginationRowModel: getPaginationRowModel(),
        onPaginationChange: setPagination,
        getSortedRowModel: getSortedRowModel(),
        onSortingChange: setSorting,
    } : {
        getSortedRowMode: getSortedRowModel(),
        onSortingChange: setSorting,
    }

    let tableState = paginated ? {
        pagination,
        sorting,
        globalFilter,
    } : {
        sorting,
        globalFilter,
    }
    const table = useReactTable({
        data: data,
        columns: columns,
        getCoreRowModel: getCoreRowModel(),
        state: tableState,
        enableSorting: sorted,
        enableGlobalFilter: filtered,
        getFilteredRowModel: getFilteredRowModel(),
        onGlobalFilterChange: setGlobalFilter,
        globalFilterFn: 'fuzzy',
        filterFns: {
            fuzzy: fuzzyFilter,
        },
        ...paginationOptions,
    })

    return (
        <div>
            <div className="px-6 py-4 flex space-x-2 items-center">
                <h3 className="flex-grow font-bold">
                    {title}
                </h3>
                {actions && (
                    <>
                        <div className="text-sm border border-slate-300 rounded py-1 px-2 focus:ring-none">
                            <DebouncedInput type="text"
                                className="focus:outline-none"
                                value={globalFilter ?? ''}
                                onChange={value => setGlobalFilter(String(value))}
                                placeholder="Search">
                            </DebouncedInput>
                            {globalFilter !== "" && (
                                <button className="button-outline button-small text-slate-300" onClick={() => setGlobalFilter("")}><i className="ri-filter-off-fill"></i></button>
                            )}
                        </div>
                    </>
                )}
                {custom_actions && custom_actions.map((elt: any) => (
                    <>
                        {elt}
                    </>
                ))}
            </div>
            <div className="relative overflow-x-auto">
                <table className="table w-full text-sm text-left rtl:text-right ">
                    <thead className="text-xs text-base-content/70 uppercase">
                        {table.getHeaderGroups().map(headerGroup => (
                            <tr key={headerGroup.id}>
                                {headerGroup.headers.map(header => (
                                    <th scope="col" className="px-6 py-3" key={header.id}>
                                        {header.isPlaceholder
                                            ? null
                                            : (
                                                <div
                                                    className={`flex ${header.column.getCanSort()
                                                        ? 'cursor-pointer select-none'
                                                        : ''
                                                        }`}
                                                    onClick={header.column.getToggleSortingHandler()}
                                                    title={
                                                        header.column.getCanSort()
                                                            ? header.column.getNextSortingOrder() === 'asc'
                                                                ? 'Sort ascending'
                                                                : header.column.getNextSortingOrder() === 'desc'
                                                                    ? 'Sort descending'
                                                                    : 'Clear sort'
                                                            : undefined
                                                    }
                                                >
                                                    <div className="flex-grow">
                                                        {
                                                            flexRender(
                                                                header.column.columnDef.header,
                                                                header.getContext()
                                                            )
                                                        }
                                                    </div>
                                                    <div>
                                                        {{
                                                            asc: <i className="ri-sort-asc"></i>,
                                                            desc: <i className="ri-sort-desc"></i>,
                                                        }[header.column.getIsSorted() as string] ?? null}
                                                    </div>
                                                </div>
                                            )
                                        }
                                    </th>
                                ))}
                            </tr>
                        ))}
                    </thead>
                    <tbody>
                        {data.length === 0 ? (
                            <tr><td colSpan={columns.length} className="text-center text-lg text-base-content/50 italic">
                                Nothing to show
                            </td></tr>
                        ) : (
                            <>
                                {
                                    table.getRowModel().rows.map(row => (
                                        <tr className="bg-white border-b border-gray-200" key={row.id}>
                                            {row.getVisibleCells().map(cell => (
                                                <td className="px-6 py-4" key={cell.id}>
                                                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                                                </td>
                                            ))}
                                        </tr>
                                    ))
                                }
                            </>
                        )}
                    </tbody>
                </table>
            </div >
            {paginated && (
                <div className="flex flex-grow justify-center items-center px-6 py-2">
                    <span className='flex items-center gap-2'>
                        Showing <b>{table.getState().pagination.pageIndex * table.getState().pagination.pageSize + 1}-{table.getState().pagination.pageIndex * table.getState().pagination.pageSize + table.getRowModel().rows.length}</b> of {table.getPrePaginationRowModel().rows.length}
                    </span>
                    <div className='flex flex-grow justify-center items-center'>
                        <PaginationButton
                            onClick={() => table.firstPage()}
                            disabled={!table.getCanPreviousPage()}>
                            <i className="ri-arrow-left-double-line"></i>
                        </PaginationButton>
                        <PaginationButton
                            onClick={() => table.previousPage()}
                            disabled={!table.getCanPreviousPage()}>
                            <i className="ri-arrow-left-s-line"></i>
                        </PaginationButton>
                        {TableUtils.getPaginationItems(table.getState().pagination.pageIndex, table.getPageCount(), 6).map((page: number) => (
                            <PaginationButton
                                key={page}
                                disabled={isNaN(page)}
                                active={table.getState().pagination.pageIndex === page}
                                onClick={() => table.setPageIndex(page)}>
                                {isNaN(page) ? '...' : page + 1}
                            </PaginationButton>
                        ))}
                        <PaginationButton
                            onClick={() => table.nextPage()}
                            disabled={!table.getCanNextPage()}>
                            <i className="ri-arrow-right-s-line"></i>
                        </PaginationButton>
                        <PaginationButton
                            onClick={() => table.lastPage()}
                            disabled={!table.getCanNextPage()}>
                            <i className="ri-arrow-right-double-line"></i>
                        </PaginationButton>
                    </div>
                    <div>
                        <select
                            value={table.getState().pagination.pageSize}
                            onChange={e => {
                                table.setPageSize(Number(e.target.value))
                            }}
                            className="form-select bg-transparent text-sm border-none">
                            {[5, 10, 25, 50, 100].map(size => (
                                <option key={size} value={size}>
                                    {size} per page
                                </option>
                            ))}
                        </select>
                    </div>
                </div>
            )
            }
        </div>
    );
}
export default Table;