import { useEffect, useState } from "react";

import PageHeader from "../components/page_header";
import Section from "../components/section";
import API from "../utils/api";
import Image from "../types/image";
import ImageList from "../components/image_list";
import Const from "../types/const";
import Utils from "../utils/utils";

function Images() {
    const [imgs, setImages] = useState<Image[]>([]);
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        function getImageList() {
            setLoading(true);
            API.images.list({ limit: 10000 }).then((response) => {
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
            <PageHeader title="Images" />
            <Section>
                <ImageList
                    actions={true}
                    data={imgs}
                    paginated={true}
                    sorted={true}
                    loading={loading}
                    deletable={true}
                    onDeleted={(uuid: string) => setImages(prev => prev.filter(img => img.uuid !== uuid))}
                />
            </Section>
        </>
    );
}

export default Images;