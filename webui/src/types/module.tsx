type Module = {
    name: string,
    module_type: string,
    backup_type: string,
    available_variants: string[],
    backup_paths: string[],
    exclude: string[],
    include: string[],
    exclude_cvs: boolean,
    variant: string,
};

export default Module;