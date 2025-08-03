import { useEffect, useState } from "react";

import Card from "../components/card";
import API from "../utils/api";
import Image from "../types/image";
import ImageList from "../components/image_list";
import Const from "../types/const";
import Utils from "../utils/utils";

function Images() {
    let [imgs, setImages] = useState<Image[]>([]);
    let [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        function getImageList() {
            setLoading(true);
            API.images.list({ limit: 10000 }).then((response: any) => {
                setImages(response.data.data ?? []);
                setLoading(false);
            }).catch(error => {
                console.log("Cannot get image list", error);
                Utils.notify(Const.CRITICAL, "Cannot get image list", error.toString())
                setImages([]);
            });
        }

        getImageList();
    }, [])

    return (
        <>
            <Card>
                <ImageList
                    title="All images"
                    actions={true}
                    data={imgs}
                    paginated={true}
                    sorted={true}
                    loading={loading}
                />
            </Card>
        </>
    );
}

export default Images;