import useIsDesktop from "./useIsDesktop";

export default function useTitleBarHeight(): string {
    return useIsDesktop() ? '45px' : '0px'
}