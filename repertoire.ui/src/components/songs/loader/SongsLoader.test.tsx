import {mantineRender} from "../../../test-utils.tsx";
import SongsLoader from "./SongsLoader.tsx";

describe('Songs Loader', () => {
  it('should render', () => {
    mantineRender(<SongsLoader />)
  })
})
