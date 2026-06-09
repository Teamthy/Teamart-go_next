import { ReactNode } from "react";
import { ChevronDown, ChevronUp, ChevronsUpDown } from "lucide-react";

interface Column<T> {
    key: keyof T;
    label: string;
    width?: string;
    render?: (value: any, row: T) => ReactNode;
    sortable?: boolean;
}

interface DataTableProps<T extends { id: string | number }> {
    columns: Column<T>[];
    data: T[];
    selectable?: boolean;
    selectedIds?: (string | number)[];
    onSelectChange?: (ids: (string | number)[]) => void;
    onSort?: (column: string, direction: "asc" | "desc") => void;
    sortColumn?: string;
    sortDirection?: "asc" | "desc";
    loading?: boolean;
    emptyMessage?: string;
}

export default function DataTable<T extends { id: string | number }>({
    columns,
    data,
    selectable = false,
    selectedIds = [],
    onSelectChange,
    onSort,
    sortColumn,
    sortDirection,
    loading = false,
    emptyMessage = "No data available",
}: DataTableProps<T>) {
    const handleSelectAll = () => {
        if (onSelectChange) {
            if (selectedIds.length === data.length) {
                onSelectChange([]);
            } else {
                onSelectChange(data.map((row) => row.id));
            }
        }
    };

    const handleSelectRow = (id: string | number) => {
        if (onSelectChange) {
            if (selectedIds.includes(id)) {
                onSelectChange(selectedIds.filter((sid) => sid !== id));
            } else {
                onSelectChange([...selectedIds, id]);
            }
        }
    };

    const handleSort = (columnKey: string, isSortable?: boolean) => {
        if (!isSortable || !onSort) return;

        if (sortColumn === columnKey) {
            onSort(columnKey, sortDirection === "asc" ? "desc" : "asc");
        } else {
            onSort(columnKey, "asc");
        }
    };

    const SortIcon = ({ columnKey, isSortable }: { columnKey: string; isSortable?: boolean }) => {
        if (!isSortable) return null;

        if (sortColumn === columnKey) {
            return sortDirection === "asc" ? (
                <ChevronUp className="w-4 h-4" />
            ) : (
                <ChevronDown className="w-4 h-4" />
            );
        }

        return <ChevronsUpDown className="w-4 h-4 opacity-40" />;
    };

    return (
        <div className="rounded-lg border border-[var(--border-light)] overflow-hidden">
            <div className="overflow-x-auto">
                <table className="w-full">
                    <thead>
                        <tr className="bg-[var(--bg-lighter)] border-b border-[var(--border-light)]">
                            {selectable && (
                                <th className="px-6 py-3 text-left">
                                    <input
                                        type="checkbox"
                                        checked={selectedIds.length === data.length && data.length > 0}
                                        onChange={handleSelectAll}
                                        className="rounded border-[var(--border-light)] cursor-pointer"
                                    />
                                </th>
                            )}
                            {columns.map((column) => (
                                <th
                                    key={String(column.key)}
                                    className={`px-6 py-3 text-left text-sm font-semibold text-[var(--text-secondary)] ${column.sortable ? "cursor-pointer hover:bg-[var(--bg-muted)]" : ""
                                        }`}
                                    style={{ width: column.width }}
                                    onClick={() =>
                                        handleSort(String(column.key), column.sortable)
                                    }
                                >
                                    <div className="flex items-center gap-2">
                                        {column.label}
                                        <SortIcon
                                            columnKey={String(column.key)}
                                            isSortable={column.sortable}
                                        />
                                    </div>
                                </th>
                            ))}
                        </tr>
                    </thead>
                    <tbody>
                        {loading ? (
                            <tr>
                                <td
                                    colSpan={columns.length + (selectable ? 1 : 0)}
                                    className="px-6 py-8 text-center"
                                >
                                    <div className="flex justify-center">
                                        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-[var(--primary)]" />
                                    </div>
                                </td>
                            </tr>
                        ) : data.length === 0 ? (
                            <tr>
                                <td
                                    colSpan={columns.length + (selectable ? 1 : 0)}
                                    className="px-6 py-8 text-center text-[var(--text-tertiary)]"
                                >
                                    {emptyMessage}
                                </td>
                            </tr>
                        ) : (
                            data.map((row) => (
                                <tr
                                    key={row.id}
                                    className="border-b border-[var(--border-light)] hover:bg-[var(--bg-lighter)] transition-colors"
                                >
                                    {selectable && (
                                        <td className="px-6 py-4">
                                            <input
                                                type="checkbox"
                                                checked={selectedIds.includes(row.id)}
                                                onChange={() => handleSelectRow(row.id)}
                                                className="rounded border-[var(--border-light)] cursor-pointer"
                                            />
                                        </td>
                                    )}
                                    {columns.map((column) => (
                                        <td
                                            key={`${row.id}-${String(column.key)}`}
                                            className="px-6 py-4 text-sm text-[var(--text-primary)]"
                                            style={{ width: column.width }}
                                        >
                                            {column.render
                                                ? column.render(row[column.key], row)
                                                : String(row[column.key])}
                                        </td>
                                    ))}
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
}
