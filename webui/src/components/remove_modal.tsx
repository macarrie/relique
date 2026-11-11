import React, { useEffect, useRef, useState } from "react";

// Generic confirmation dialog for destructive actions.
// onConfirm must return a promise: the dialog closes when it resolves and stays open if it rejects.
function RemoveModal({
    open = false,
    title = "",
    children = null as React.ReactNode,
    confirmLabel = "Delete",
    confirmDisabled = false,
    onConfirm = () => Promise.resolve() as Promise<unknown>,
    onClose = () => { },
}) {
    const dialog = useRef<HTMLDialogElement>(null);
    const [removing, setRemoving] = useState<boolean>(false);

    useEffect(() => {
        if (open) {
            dialog.current?.showModal();
        } else {
            dialog.current?.close();
        }
    }, [open]);

    function confirm() {
        setRemoving(true);
        onConfirm().then(() => {
            onClose();
        }).catch(() => {
            // Error is notified by the caller, keep the dialog open
        }).finally(() => {
            setRemoving(false);
        });
    }

    return (
        <dialog ref={dialog} className="modal" onClose={() => { if (open) onClose(); }}>
            <div className="modal-box max-w-2xl">
                <h3 className="font-bold text-lg">{title}</h3>
                <div className="py-4 space-y-4">
                    {children}
                    <p className="font-semibold text-error">This action is not recoverable.</p>
                </div>
                <div className="modal-action">
                    <button type="button" className="btn" disabled={removing} onClick={onClose}>Cancel</button>
                    <button type="button" className="btn btn-error" disabled={removing || confirmDisabled} onClick={confirm}>
                        {removing && <span className="loading loading-spinner loading-sm"></span>}
                        <i className="ri-delete-bin-line"></i>
                        {confirmLabel}
                    </button>
                </div>
            </div>
        </dialog>
    );
}

export default RemoveModal;
