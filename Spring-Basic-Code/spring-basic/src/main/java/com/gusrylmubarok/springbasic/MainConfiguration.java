package com.gusrylmubarok.springbasic;

import com.gusrylmubarok.springbasic.configuration.BarConfiguration;
import com.gusrylmubarok.springbasic.configuration.FooConfiguration;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Import;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Configuration
@Import({
        FooConfiguration.class,
        BarConfiguration.class
})
public class MainConfiguration {
}