package com.gusrylmubarok.springbasic.bean.scopes;

import com.gusrylmubarok.springbasic.data.Foo;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Scope;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Slf4j
@Configuration
public class ScopeConfiguration {

    @Bean
    @Scope("protorype")
    public Foo foo() {
        log.info("Create new Foo");
        return new Foo();
    }

}
