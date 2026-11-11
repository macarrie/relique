import { useEffect, useState } from "react";
import { AxiosError, AxiosResponse } from "axios";

import Job from "../types/job";
import Image from "../types/image";
import Const from "../types/const";
import API from "../utils/api";
import Utils from "../utils/utils";
import RemoveModal from "./remove_modal";
import TextPlaceholder from "./text_placeholder";
import PropertyList from "./property_list";

// Asks confirmation before removing an image, with a recap of what is removed and what is kept.
// onRemoved receives the image uuid and whether the job that generated it was removed too.
type Props = {
    image: Image | null,
    onClose: () => void,
    onRemoved: (uuid: string, jobRemoved: boolean) => void,
};

function ImageRemoveModal({ image, onClose, onRemoved }: Props) {
    const [linkedJob, setLinkedJob] = useState<Job | null>(null);
    const [linkedLoading, setLinkedLoading] = useState<boolean>(image !== null);
    const [cleanJob, setCleanJob] = useState<boolean>(false);
    const [prevImage, setPrevImage] = useState<Image | null>(image);

    // Reset the dialog when another image is selected
    if (image !== prevImage) {
        setPrevImage(image);
        setCleanJob(false);
        setLinkedJob(null);
        setLinkedLoading(image !== null);
    }

    useEffect(() => {
        if (!image) {
            return;
        }

        // Images share their uuid with the job that generated them
        API.jobs.get(image.uuid).then((response: AxiosResponse<Job>) => {
            setLinkedJob(response.data);
        }).catch(() => {
            setLinkedJob(null);
        }).finally(() => {
            setLinkedLoading(false);
        });
    }, [image]);

    function remove(): Promise<unknown> {
        if (!image) {
            return Promise.resolve();
        }

        const jobRemoved = cleanJob && linkedJob !== null;
        return API.images.delete(image.uuid, jobRemoved).then(() => {
            Utils.notify(Const.OK, "Image removed", jobRemoved ? "Image and its job have been removed" : "Image has been removed");
            onRemoved(image.uuid, jobRemoved);
        }).catch((error: AxiosError<{ error?: string }>) => {
            console.log("Cannot remove image", error);
            Utils.notify(Const.CRITICAL, "Cannot remove image", error.response?.data?.error ?? error.toString());
            throw error;
        });
    }

    return (
        <RemoveModal
            open={image !== null}
            title="Image removal"
            confirmLabel={cleanJob && linkedJob ? "Delete image and job" : "Delete image"}
            confirmDisabled={linkedLoading}
            onConfirm={remove}
            onClose={onClose}
        >
            <p>
                The following image will be removed. Its backup data (<b>{Utils.formatSize(image?.size_on_disk ?? 0)}</b>) will
                be permanently deleted from repository <b>{image?.repository?.name}</b>:
            </p>
            <PropertyList items={[
                { label: "ID", value: image?.uuid, mono: true },
                { label: "Client", value: image?.client?.name },
                { label: "Module", value: image?.module?.name },
                { label: "Date", value: image && <span className="tabular-nums">{Utils.formatDate(image.created_at)}</span> },
                { label: "Size", value: <span className="tabular-nums">{Utils.formatSize(image?.size_on_disk ?? 0)}</span> },
            ]} />

            {linkedLoading && <TextPlaceholder />}
            {!linkedLoading && linkedJob && (
                <div className="space-y-2">
                    <label className="label cursor-pointer text-base-content">
                        <input type="checkbox" className="checkbox checkbox-error checkbox-sm" checked={cleanJob} onChange={e => setCleanJob(e.target.checked)} />
                        Also remove the job that generated this image
                    </label>
                    {!cleanJob && (
                        <p className="text-muted">
                            The job is kept in the job history, with its logs and statistics.
                        </p>
                    )}
                </div>
            )}
        </RemoveModal>
    );
}

export default ImageRemoveModal;
