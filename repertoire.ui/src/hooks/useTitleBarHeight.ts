import useIsDesktop from './useIsDesktop'

export default function useTitleBarHeight() {
  return useIsDesktop() ? '45px' : '0px'
}
