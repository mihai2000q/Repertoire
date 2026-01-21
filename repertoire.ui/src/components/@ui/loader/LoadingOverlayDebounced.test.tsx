import { mantineRender } from '../../../test-utils.tsx'
import LoadingOverlayDebounced from './LoadingOverlayDebounced.tsx'

describe('Loading Overlay Debounced', () => {
  it('should render', () => {
    mantineRender(<LoadingOverlayDebounced />)
  })
})
