import Module from "../types/module";
import PropertyList, { Property } from "./property_list";
import YesNo from "./yes_no";

function pathList(paths?: string[]) {
    if (!paths || paths.length === 0) {
        return undefined;
    }

    return (
        <ul>
            {paths.map((path: string) => <li key={path}>{path}</li>)}
        </ul>
    );
}

// installed=true describes a module installed on the server (shows available variants);
// otherwise it describes a module configured on a client (shows the selected variant).
function ModuleProperties({
    module,
    loading = false,
    installed = false,
}: {
    module: Module,
    loading?: boolean,
    installed?: boolean,
}) {
    const items: Property[] = [
        { label: "Module type", value: module?.module_type },
        installed
            ? { label: "Variants", value: (module?.available_variants ?? []).join(", ") }
            : { label: "Variant", value: module?.variant || "default" },
        { label: "Backup type", value: module?.backup_type },
        { label: "Backup paths", value: pathList(module?.backup_paths), mono: true },
        { label: "Exclusions", value: pathList(module?.exclude), mono: true },
        { label: "Inclusions", value: pathList(module?.include), mono: true },
        { label: "Exclude CVS", value: <YesNo value={module?.exclude_cvs} /> },
    ];

    return <PropertyList items={items} loading={loading} />;
}

export default ModuleProperties;
