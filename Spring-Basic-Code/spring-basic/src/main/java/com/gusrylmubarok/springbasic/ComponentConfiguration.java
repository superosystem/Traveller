package com.gusrylmubarok.springbasic;

import com.gusrylmubarok.springbasic.data.MultiFoo;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Import;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Configuration
@ComponentScan(basePackages = {
        "com.gusrylmubarok.springbasic.repository",
        "com.gusrylmubarok.springbasic.service",
        "com.gusrylmubarok.springbasic.configuration",
})
@Import(MultiFoo.class)
public class ComponentConfiguration {
}
