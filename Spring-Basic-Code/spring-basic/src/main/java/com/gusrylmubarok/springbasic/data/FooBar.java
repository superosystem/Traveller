package com.gusrylmubarok.springbasic.data;

import lombok.AllArgsConstructor;
import lombok.Data;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Data
@AllArgsConstructor
public class FooBar {
    private Foo foo;

    private Bar bar;
}
