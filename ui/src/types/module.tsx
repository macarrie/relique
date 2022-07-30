type Module = {
    module_type :string,
    name :string,
    backup_type :string,
    schedules :any[],
    available_variants :any[],
    backup_paths :string[],
    pre_backup_script :string,
    post_backup_script :string,
    pre_restore_script :string,
    post_restore_script :string,
    variant :string,
    params :any,
};

export default Module;
