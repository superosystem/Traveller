package com.gusrylmubarok.springbasic;

import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Configuration
@ComponentScan(basePackages = {
        "com.gusrylmubarok.springbasic.configuration"
})
public class ScanConfiguration {
}
