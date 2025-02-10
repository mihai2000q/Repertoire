import {renderHook} from "@testing-library/react";
import useInstrumentIcon from "./useInstrumentIcon.tsx";

describe('use Instrument Icon', () => {
  it('should return the icon of the instrument', () => {
    const { result } = renderHook(() => useInstrumentIcon())

    // return specific icon
    expect(result.current('Voice')).not.toBeUndefined()
    expect(result.current('Voice')).not.toBeNull()

    // return default icon
    expect(result.current('random')).not.toBeUndefined()
    expect(result.current('random')).not.toBeNull()
  })
})
