import { useEffect, useState } from "react";

import Card from "../components/card";
import API from "../utils/api";
import Image from "../types/image";
import ImageList from "../components/image_list";

function Images() {
    let [imgs, setImages] = useState<Image[]>([]);

    useEffect(() => {
        function getImageList() {
            API.images.list({ limit: 10000 }).then((response: any) => {
                setImages(response.data.data ?? []);
            }).catch(error => {
                console.log("Cannot get job list", error);
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
                />
            </Card>
        </>
    );
}

export default Images;