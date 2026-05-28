type Column<T> = {
    header: string;
    accessor: keyof T | ((row: T) => React.ReactNode);
    className?: string;
};

interface DataTableProps<T extends Record<string, unknown>> {
    columns: Column<T>[];
    rows: T[];
}

export default function DataTable<T extends Record<string, unknown>>({ columns, rows }: DataTableProps<T>) {
    return (
        <div className="overflow-hidden rounded-[24px] border border-slate-200 bg-white">
            <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-slate-200 text-left text-sm">
                    <thead className="bg-slate-50">
                        <tr>
                            {columns.map((column, index) => (
                                <th key={`${String(column.header)}-${index}`} className="px-4 py-3 font-semibold text-slate-700">
                                    {column.header}
                                </th>
                            ))}
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-100 bg-white">
                        {rows.map((row, rowIndex) => (
                            <tr key={rowIndex} className="hover:bg-slate-50">
                                {columns.map((column, columnIndex) => {
                                    const value = typeof column.accessor === "function" ? column.accessor(row) : row[column.accessor];
                                    return (
                                        <td key={`${rowIndex}-${columnIndex}`} className={`px-4 py-4 ${column.className ?? "text-slate-700"}`}>
                                            {value as React.ReactNode}
                                        </td>
                                    );
                                })}
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
}
