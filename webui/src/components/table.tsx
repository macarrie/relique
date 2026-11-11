import { useState } from "react";
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

import TableUtils, { TableColumns } from "../utils/table";
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
const fuzzyFilter: FilterFn<unknown> = (row, columnId, value, addMeta) => {
    // Rank the item
    const itemRank = rankItem(row.getValue(columnId), value)

    // Store the itemRank info
    addMeta({
        itemRank,
    })

    // Return if the item should be filtered in/out
    return itemRank.passed
}

type TableProps<T> = {
    actions?: boolean,
    filtered?: boolean,
    paginated?: boolean,
    sorted?: boolean,
    columns: TableColumns<T>,
    data: T[],
    defaultPageSize?: number,
};

function Table<T>({
    actions = true,
    filtered = true,
    paginated = true,
    sorted = true,
    columns,
    data,
    defaultPageSize = Const.DEFAULT_PAGE_SIZE,
}: TableProps<T>) {
    const [pagination, setPagination] = useState<PaginationState>({
        pageIndex: 0,
        pageSize: defaultPageSize,
    })
    const [sorting, setSorting] = useState<SortingState>([]);
    const [globalFilter, setGlobalFilter] = useState('');

    const paginationOptions = paginated ? {
        getPaginationRowModel: getPaginationRowModel(),
        onPaginationChange: setPagination,
        getSortedRowModel: getSortedRowModel(),
        onSortingChange: setSorting,
    } : {
        getSortedRowModel: getSortedRowModel(),
        onSortingChange: setSorting,
    }

    const tableState = paginated ? {
        pagination,
        sorting,
        globalFilter,
    } : {
        sorting,
        globalFilter,
    }
    // TanStack Table returns non-memoizable functions, so React Compiler skips this
    // component. We don't use the compiler, so this is expected.
    // eslint-disable-next-line react-hooks/incompatible-library
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

    const { pageIndex, pageSize } = table.getState().pagination ?? { pageIndex: 0, pageSize: 0 };
    const totalRows = table.getPrePaginationRowModel().rows.length;
    const firstRow = totalRows === 0 ? 0 : pageIndex * pageSize + 1;
    const lastRow = pageIndex * pageSize + table.getRowModel().rows.length;

    return (
        <div className="space-y-2">
            {actions && (
                <label className="input input-sm w-64 max-w-full">
                    <i className="ri-search-line text-muted"></i>
                    <DebouncedInput type="search"
                        className="grow"
                        value={globalFilter ?? ''}
                        onChange={value => setGlobalFilter(String(value))}
                        placeholder="Filter">
                    </DebouncedInput>
                </label>
            )}
            <div className="overflow-x-auto">
                <table className="w-full text-left">
                    <thead>
                        {table.getHeaderGroups().map(headerGroup => (
                            <tr key={headerGroup.id} className="border-b border-base-300">
                                {headerGroup.headers.map(header => (
                                    <th scope="col" className="px-3 py-2 text-label whitespace-nowrap" key={header.id}>
                                        {header.isPlaceholder
                                            ? null
                                            : (
                                                <div
                                                    className={`flex items-center gap-1 ${header.column.getCanSort()
                                                        ? 'cursor-pointer select-none hover:text-base-content'
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
                                                    {
                                                        flexRender(
                                                            header.column.columnDef.header,
                                                            header.getContext()
                                                        )
                                                    }
                                                    {{
                                                        asc: <i className="ri-arrow-up-s-line text-sm"></i>,
                                                        desc: <i className="ri-arrow-down-s-line text-sm"></i>,
                                                    }[header.column.getIsSorted() as string] ?? null}
                                                </div>
                                            )
                                        }
                                    </th>
                                ))}
                            </tr>
                        ))}
                    </thead>
                    <tbody>
                        {table.getRowModel().rows.length === 0 ? (
                            <tr><td colSpan={columns.length} className="px-3 py-6 text-center text-muted">
                                {globalFilter !== "" ? "No match for this filter" : "Nothing to show"}
                            </td></tr>
                        ) : (
                            table.getRowModel().rows.map(row => (
                                <tr className="border-b border-base-300 hover:bg-base-content/[0.04]" key={row.id}>
                                    {row.getVisibleCells().map(cell => (
                                        <td className="px-3 py-2 align-middle whitespace-nowrap" key={cell.id}>
                                            {flexRender(cell.column.columnDef.cell, cell.getContext())}
                                        </td>
                                    ))}
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>
            {paginated && (
                <div className="flex flex-wrap items-center gap-4 text-muted">
                    <span className="tabular-nums">
                        <span className="text-base-content font-medium">{firstRow}–{lastRow}</span> of {totalRows}
                    </span>
                    {table.getPageCount() > 1 && (
                        <div className="flex flex-grow justify-center items-center gap-0.5">
                            <button className="btn btn-ghost btn-sm btn-square" aria-label="First page"
                                onClick={() => table.firstPage()}
                                disabled={!table.getCanPreviousPage()}>
                                <i className="ri-arrow-left-double-line"></i>
                            </button>
                            <button className="btn btn-ghost btn-sm btn-square" aria-label="Previous page"
                                onClick={() => table.previousPage()}
                                disabled={!table.getCanPreviousPage()}>
                                <i className="ri-arrow-left-s-line"></i>
                            </button>
                            {TableUtils.getPaginationItems(pageIndex, table.getPageCount(), 6).map((page: number, i: number) => (
                                isNaN(page) ? (
                                    <span key={`gap-${i}`} className="w-8 text-center">…</span>
                                ) : (
                                    <button key={page}
                                        className={`btn btn-ghost btn-sm btn-square font-normal tabular-nums ${pageIndex === page ? "btn-active font-semibold text-base-content" : ""}`}
                                        aria-current={pageIndex === page ? "page" : undefined}
                                        onClick={() => table.setPageIndex(page)}>
                                        {page + 1}
                                    </button>
                                )
                            ))}
                            <button className="btn btn-ghost btn-sm btn-square" aria-label="Next page"
                                onClick={() => table.nextPage()}
                                disabled={!table.getCanNextPage()}>
                                <i className="ri-arrow-right-s-line"></i>
                            </button>
                            <button className="btn btn-ghost btn-sm btn-square" aria-label="Last page"
                                onClick={() => table.lastPage()}
                                disabled={!table.getCanNextPage()}>
                                <i className="ri-arrow-right-double-line"></i>
                            </button>
                        </div>
                    )}
                    <select
                        value={pageSize}
                        onChange={e => {
                            table.setPageSize(Number(e.target.value))
                        }}
                        className="select select-sm select-ghost w-auto ml-auto">
                        {[5, 10, 25, 50, 100].map(size => (
                            <option key={size} value={size}>
                                {size} per page
                            </option>
                        ))}
                    </select>
                </div>
            )}
        </div>
    );
}
export default Table;
