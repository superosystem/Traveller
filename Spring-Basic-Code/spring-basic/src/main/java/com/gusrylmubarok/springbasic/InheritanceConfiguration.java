package com.gusrylmubarok.springbasic;

import com.gusrylmubarok.springbasic.service.MerchantServiceImpl;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Import;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Configuration
@Import(MerchantServiceImpl.class)
public class InheritanceConfiguration {

}
