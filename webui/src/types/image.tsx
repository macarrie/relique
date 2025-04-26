import Client from "./client";
import Module from "./module";
import Repository from "./repository";

type Image = {
    uuid: string,
    client: Client,
    module: Module,
    repository: Repository,
    number_of_elements: string,
    number_of_files: boolean,
    number_of_folders: string,
    size_on_disk: number,
    created_at: any,
};

export default Image;