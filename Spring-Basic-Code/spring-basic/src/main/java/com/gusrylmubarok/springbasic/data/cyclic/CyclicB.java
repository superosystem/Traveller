package com.gusrylmubarok.springbasic.data.cyclic;

import lombok.AllArgsConstructor;
import lombok.Data;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Data
@AllArgsConstructor
public class CyclicB {
    private CyclicC cyclicC;
}
