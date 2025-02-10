import {renderHook} from "@testing-library/react";
import useInstrumentIcon from "./useInstrumentIcon.tsx";
import {Instrument} from "../types/models/Song.ts";

describe('use Instrument Icon', () => {
  it('should return the icon of the instrument', () => {
    const { result } = renderHook(() => useInstrumentIcon())

    // return specific icon
    expect(result.current('Voice')).not.toBeUndefined()
    expect(result.current('Voice')).not.toBeNull()

    // return specific icon with instrument type
    const instrument: Instrument = {
      id: '1',
      name: 'Electric Guitar'
    }
    expect(result.current(instrument)).not.toBeUndefined()
    expect(result.current(instrument)).not.toBeNull()

    // return default icon
    expect(result.current('random')).not.toBeUndefined()
    expect(result.current('random')).not.toBeNull()
  })
})
