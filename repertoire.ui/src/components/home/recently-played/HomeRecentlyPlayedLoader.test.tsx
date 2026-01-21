import { mantineRender } from '../../../test-utils.tsx'
import HomeRecentlyPlayedLoader from './HomeRecentlyPlayedLoader.tsx'

describe('Home Recently Played Loader', () => {
  it('should render', () => {
    mantineRender(<HomeRecentlyPlayedLoader />)
  })
})
