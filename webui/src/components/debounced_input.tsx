import { useEffect, useEffectEvent, useState } from "react"

function DebouncedInput({
    value: initialValue,
    onChange,
    debounce = 500,
    ...props
}: {
    value: string | number
    onChange: (value: string | number) => void
    debounce?: number
} & Omit<React.InputHTMLAttributes<HTMLInputElement>, 'onChange'>) {
    const [value, setValue] = useState(initialValue)
    const [prevInitialValue, setPrevInitialValue] = useState(initialValue)

    // Reset the local value when the value prop changes from the outside
    if (initialValue !== prevInitialValue) {
        setPrevInitialValue(initialValue)
        setValue(initialValue)
    }

    // Read the latest onChange without restarting the debounce timer when it changes
    const emitChange = useEffectEvent((v: string | number) => onChange(v))

    useEffect(() => {
        const timeout = setTimeout(() => {
            emitChange(value)
        }, debounce)

        return () => clearTimeout(timeout)
    }, [value, debounce])

    return (
        <input {...props} value={value} onChange={e => setValue(e.target.value)} />
    )
}

export default DebouncedInput;
